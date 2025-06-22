package battlebots

import (
	"context"

	"github.com/z5labs/battlebots/sdk/battlebots-go/battlebotspb"
)

type DefaultState struct{}

func (state *DefaultState) Init(ctx context.Context) error {
	return nil
}

func (state *DefaultState) ApplyEvent(ctx context.Context, ev *battlebotspb.StateChangeEvent) error {
	return nil
}
