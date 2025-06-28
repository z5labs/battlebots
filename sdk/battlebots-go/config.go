package battlebots

import (
	"github.com/z5labs/humus/job"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Config
type Config struct {
	job.Config `config:",squash"`

	BattleService struct {
		Endpoint string `config:"endpoint"`
	} `config:"battle_service"`
}

func (cfg Config) battleServiceClientConn() (*grpc.ClientConn, error) {
	return grpc.NewClient(
		cfg.BattleService.Endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler(
			otelgrpc.WithMessageEvents(otelgrpc.ReceivedEvents, otelgrpc.SentEvents),
		)),
	)
}
