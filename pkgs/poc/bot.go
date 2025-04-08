package poc

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/z5labs/battlebots/pkgs/battlebotspb"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"golang.org/x/sync/errgroup"
)

type BotOption func(*Bot)

func MinWaitFor(d time.Duration) BotOption {
	return func(b *Bot) {
		b.minWaitFor = d
	}
}

func MaxWaitFor(d time.Duration) BotOption {
	return func(b *Bot) {
		b.maxWaitFor = d
	}
}

type Bot struct {
	battlebotspb.BattleClient

	log        *slog.Logger
	rand       *rand.Rand
	minWaitFor time.Duration
	maxWaitFor time.Duration
}

func NewBot(client battlebotspb.BattleClient, randSrc rand.Source, opts ...BotOption) *Bot {
	b := &Bot{
		BattleClient: client,
		log:          otelslog.NewLogger("poc"),
		rand:         rand.New(randSrc),
		minWaitFor:   16 * time.Millisecond,
		maxWaitFor:   100 * time.Millisecond,
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func (b *Bot) Run(ctx context.Context) error {
	eg, egctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		var subscription battlebotspb.StateChangeSubscription
		events, err := b.State(egctx, &subscription)
		if err != nil {
			return err
		}

		for {
			event, err := events.Recv()
			if err != nil {
				return err
			}

			b.processEvent(egctx, event)
		}
	})

	eg.Go(func() error {
		for {
			waitFor := time.Duration(randInt64Within(b.rand, int64(b.minWaitFor), int64(b.maxWaitFor))) * time.Millisecond
			select {
			case <-egctx.Done():
				return egctx.Err()
			case <-time.After(waitFor):
			}

			err := b.randomlyMove(egctx)
			if err != nil {
				return err
			}
		}
	})

	return eg.Wait()
}

func randInt64Within(r *rand.Rand, min, max int64) int64 {
	return r.Int64N(max-min) + min
}

func (b *Bot) processEvent(ctx context.Context, event *battlebotspb.StateChangeEvent) {
	spanCtx, span := otel.Tracer("cmd").Start(ctx, "Bot.processEvent")
	defer span.End()

	if !event.HasPosition() {
		b.log.WarnContext(spanCtx, "expected to only receive position state changes")
		return
	}

	bot := event.GetBot()
	pos := event.GetPosition()

	b.log.InfoContext(
		spanCtx,
		"received position state change event",
		slog.String("bot.id", bot.GetId()),
		slog.Float64("position.x0", pos.GetX0()),
		slog.Float64("position.x1", pos.GetX1()),
	)
}

func (b *Bot) randomlyMove(ctx context.Context) error {
	spanCtx, span := otel.Tracer("cmd").Start(ctx, "Bot.randomlyMove")
	defer span.End()

	var velocity battlebotspb.Vector
	velocity.SetX0(rand.Float64())
	velocity.SetX1(rand.Float64())

	var req battlebotspb.MoveRequest
	req.SetVelocity(&velocity)

	_, err := b.Move(spanCtx, &req)
	return err
}
