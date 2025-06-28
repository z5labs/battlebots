package battle

import (
	"context"
	"log/slog"

	"github.com/z5labs/battlebots/sdk/battlebots-go/battlebotspb"
	"github.com/z5labs/humus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// Service
type Service struct {
	battlebotspb.UnimplementedBattleServer

	tracer trace.Tracer
	log    *slog.Logger
}

// NewService
func NewService() *Service {
	return &Service{
		tracer: otel.Tracer("battle"),
		log:    humus.Logger("battle"),
	}
}

// Move implements [battlebotspb.BattleServer].
func (s *Service) Move(context.Context, *battlebotspb.MoveRequest) (*battlebotspb.MoveResponse, error) {
	panic("unimplemented")
}

// State implements [battlebotspb.BattleServer].
func (s *Service) State(*battlebotspb.StateChangeSubscription, grpc.ServerStreamingServer[battlebotspb.StateChangeEvent]) error {
	panic("unimplemented")
}
