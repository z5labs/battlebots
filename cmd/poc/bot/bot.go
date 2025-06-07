package bot

import (
	"context"
	crand "crypto/rand"
	"io"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/z5labs/battlebots/pkgs/battlebotspb"
	"github.com/z5labs/humus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
)

type Option func(*Handler)

func RandSource(src rand.Source) Option {
	return func(h *Handler) {
		h.rand = rand.New(src)
	}
}

func MinWaitFor(d time.Duration) Option {
	return func(h *Handler) {
		h.minWaitFor = d
	}
}

func MaxWaitFor(d time.Duration) Option {
	return func(h *Handler) {
		h.maxWaitFor = d
	}
}

// Handler
type Handler struct {
	battlebotspb.BattleClient

	tracer     trace.Tracer
	log        *slog.Logger
	rand       *rand.Rand
	minWaitFor time.Duration
	maxWaitFor time.Duration
}

func NewHandler(client battlebotspb.BattleClient, opts ...Option) (*Handler, error) {
	var seed [32]byte
	_, err := io.ReadFull(crand.Reader, seed[:])
	if err != nil {
		return nil, err
	}

	h := &Handler{
		BattleClient: client,
		tracer:       otel.Tracer("bot"),
		log:          humus.Logger("bot"),
		rand:         rand.New(rand.NewChaCha8(seed)),
		minWaitFor:   16 * time.Millisecond,
		maxWaitFor:   100 * time.Millisecond,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h, nil
}

// Handle implements [job.Handler].
func (h *Handler) Handle(ctx context.Context) error {
	eg, egctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		var subscription battlebotspb.StateChangeSubscription
		events, err := h.State(egctx, &subscription)
		if err != nil {
			return err
		}

		for {
			event, err := events.Recv()
			if err != nil {
				return err
			}

			h.processEvent(egctx, event)
		}
	})

	eg.Go(func() error {
		for {
			waitFor := time.Duration(randInt64Within(h.rand, int64(h.minWaitFor), int64(h.maxWaitFor))) * time.Millisecond
			select {
			case <-egctx.Done():
				return egctx.Err()
			case <-time.After(waitFor):
			}

			err := h.randomlyMove(egctx)
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

func (h *Handler) processEvent(ctx context.Context, event *battlebotspb.StateChangeEvent) {
	spanCtx, span := h.tracer.Start(ctx, "Handler.processEvent")
	defer span.End()

	if !event.HasPosition() {
		h.log.WarnContext(spanCtx, "expected to only receive position state changes")
		return
	}

	bot := event.GetBot()
	pos := event.GetPosition()

	h.log.InfoContext(
		spanCtx,
		"received position state change event",
		slog.String("bot.id", bot.GetId()),
		slog.Float64("position.x0", pos.GetX0()),
		slog.Float64("position.x1", pos.GetX1()),
	)
}

func (h *Handler) randomlyMove(ctx context.Context) error {
	spanCtx, span := h.tracer.Start(ctx, "Handler.randomlyMove")
	defer span.End()

	var velocity battlebotspb.Vector
	velocity.SetX0(rand.Float64())
	velocity.SetX1(rand.Float64())

	var req battlebotspb.MoveRequest
	req.SetVelocity(&velocity)

	_, err := h.Move(spanCtx, &req)
	return err
}
