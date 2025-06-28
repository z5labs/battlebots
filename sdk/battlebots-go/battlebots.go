package battlebots

import (
	"context"
	"io"
	"log/slog"

	"github.com/z5labs/battlebots/sdk/battlebots-go/battlebotspb"
	"github.com/z5labs/humus"

	"github.com/z5labs/bedrock/lifecycle"
	"github.com/z5labs/humus/job"
	"github.com/z5labs/sdk-go/concurrent"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// Configer
type Configer interface {
	job.Configer

	battleServiceClientConn() (*grpc.ClientConn, error)
}

// State
type State interface {
	Init(context.Context) error

	ApplyEvent(context.Context, *battlebotspb.StateChangeEvent) error
}

// Controller
type Controller[S State] interface {
	Tick(context.Context, S) error
}

// Run
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
			log:    humus.Logger("battlebots"),
			tracer: otel.Tracer("battlebots"),
			client: bc,
			ctlr:   ctlr,
		}

		return job.NewApp(h), nil
	})
}

type handler[S State] struct {
	log    *slog.Logger
	tracer trace.Tracer
	client battlebotspb.BattleClient
	ctlr   Controller[S]
}

func (h *handler[S]) Handle(ctx context.Context) error {
	var subscription battlebotspb.StateChangeSubscription
	stream, err := h.client.State(ctx, &subscription)
	if err != nil {
		h.log.ErrorContext(ctx, "failed to subscribe to state change events", slog.Any("error", err))
		return err
	}

	eventCh := make(chan *battlebotspb.StateChangeEvent)

	var lg concurrent.LazyGroup
	lg.Go(func(ctx context.Context) error {
		defer close(eventCh)

		for {
			err := h.receiveEvent(ctx, stream, eventCh)
			if err != nil {
				return err
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
			err := h.tick(ctx, state, eventCh)
			if err != nil {
				return err
			}
		}
	})

	return lg.Wait(ctx)
}

func (h *handler[S]) receiveEvent(ctx context.Context, stream grpc.ServerStreamingClient[battlebotspb.StateChangeEvent], eventCh chan<- *battlebotspb.StateChangeEvent) error {
	spanCtx, span := h.tracer.Start(ctx, "handler.receiveEvent")
	defer span.End()

	event, err := stream.Recv()
	if err != nil {
		span.RecordError(err)
		return err
	}

	select {
	case <-spanCtx.Done():
		return spanCtx.Err()
	case eventCh <- event:
		return nil
	}
}

func (h *handler[S]) tick(ctx context.Context, state S, eventCh <-chan *battlebotspb.StateChangeEvent) error {
	spanCtx, span := h.tracer.Start(ctx, "handler.tick")
	defer span.End()

	event, done := pollEvent(spanCtx, eventCh)
	if done {
		return nil
	}
	if event != nil {
		err := state.ApplyEvent(spanCtx, event)
		if err != nil {
			span.RecordError(err)
			h.log.ErrorContext(spanCtx, "failed to apply event", slog.Any("error", err))
			return err
		}
	}

	err := h.ctlr.Tick(spanCtx, state)
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
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
