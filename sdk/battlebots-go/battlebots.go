package battlebots

import (
	"context"
	"io"

	"github.com/z5labs/battlebots/sdk/battlebots-go/battlebotspb"
	"github.com/z5labs/bedrock/lifecycle"
	"google.golang.org/grpc"

	"github.com/z5labs/humus/job"
	"github.com/z5labs/sdk-go/concurrent"
	"go.opentelemetry.io/otel"
)

type Configer interface {
	job.Configer

	battleServiceClientConn() (*grpc.ClientConn, error)
}

type State interface {
	Init(context.Context) error

	ApplyEvent(context.Context, *battlebotspb.StateChangeEvent) error
}

type Controller[S State] interface {
	Tick(context.Context, S) error
}

func Run[C Configer, S State, T Controller[S]](r io.Reader, f func(context.Context, C, *Bot) (T, error)) {
	job.Run(r, func(ctx context.Context, cfg C) (*job.App, error) {
		cc, err := cfg.battleServiceClientConn()
		if err != nil {
			return nil, err
		}
		lifecycle.TryCloseOnPostRun(ctx, cc)

		bc := battlebotspb.NewBattleClient(cc)

		bot := &Bot{
			tracer: otel.Tracer("battlebots"),
			client: bc,
		}

		ctlr, err := f(ctx, cfg, bot)
		if err != nil {
			return nil, err
		}

		h := &handler[S]{
			client: bc,
			ctlr:   ctlr,
		}

		return job.NewApp(h), nil
	})
}

type handler[S State] struct {
	client battlebotspb.BattleClient
	ctlr   Controller[S]
}

func (h *handler[S]) Handle(ctx context.Context) error {
	var subscription battlebotspb.StateChangeSubscription
	events, err := h.client.State(ctx, &subscription)
	if err != nil {
		return err
	}

	eventCh := make(chan *battlebotspb.StateChangeEvent)

	var lg concurrent.LazyGroup
	lg.Go(func(ctx context.Context) error {
		defer close(eventCh)

		for {
			event, err := events.Recv()
			if err != nil {
				return err
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case eventCh <- event:
			}
		}
	})

	lg.Go(func(ctx context.Context) error {
		var state S
		err := state.Init(ctx)
		if err != nil {
			return err
		}

		for {
			event, done := pollEvent(ctx, eventCh)
			if done {
				return nil
			}
			if event != nil {
				// TODO: apply event to state
			}

			err := h.ctlr.Tick(ctx, state)
			if err != nil {
				return err
			}
		}
	})

	return lg.Wait(ctx)
}

func pollEvent(ctx context.Context, eventCh <-chan *battlebotspb.StateChangeEvent) (*battlebotspb.StateChangeEvent, bool) {
	select {
	case <-ctx.Done():
	case event, ok := <-eventCh:
		return event, !ok
	default:
	}

	return nil, false
}
