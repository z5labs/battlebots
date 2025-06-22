package battlebots

import (
	"context"

	"github.com/z5labs/battlebots/sdk/battlebots-go/battlebotspb"

	"go.opentelemetry.io/otel/trace"
)

type Bot struct {
	tracer trace.Tracer
	client battlebotspb.BattleClient
}

func (b *Bot) Move(ctx context.Context, dx, dy float64) error {
	spanCtx, span := b.tracer.Start(ctx, "Bot.Move")
	defer span.End()

	var velocity battlebotspb.Vector
	velocity.SetX0(dx)
	velocity.SetX1(dy)

	var req battlebotspb.MoveRequest
	req.SetVelocity(&velocity)

	_, err := b.client.Move(spanCtx, &req)
	return err
}
