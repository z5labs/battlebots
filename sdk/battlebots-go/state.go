package battlebots

import (
	"context"

	"github.com/z5labs/battlebots/sdk/battlebots-go/battlebotspb"
)

// Entity
type Entity struct {
	Position *battlebotspb.Vector
	Velocity *battlebotspb.Vector
}

// DefaultState
type DefaultState struct {
	Entities map[string]*Entity
}

// Init implements the [State] interface.
func (state *DefaultState) Init(ctx context.Context) error {
	state.Entities = make(map[string]*Entity)
	return nil
}

// ApplyEvent implements the [State] interface.
func (state *DefaultState) ApplyEvent(ctx context.Context, ev *battlebotspb.StateChangeEvent) error {
	switch ev.WhichEvent() {
	case battlebotspb.StateChangeEvent_Position_case:
		bot := ev.GetBot()

		ent, ok := state.Entities[bot.GetId()]
		if !ok {
			ent = &Entity{}
			state.Entities[bot.GetId()] = ent
		}

		ent.Position = ev.GetPosition()
	}

	return nil
}
