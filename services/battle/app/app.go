package app

import (
	"context"

	"github.com/z5labs/battlebots/sdk/battlebots-go/battlebotspb"
	"github.com/z5labs/battlebots/services/battle/battle"
	"github.com/z5labs/humus/grpc"
)

type Config struct {
	grpc.Config `config:",squash"`
}

func Init(ctx context.Context, cfg *Config) (*grpc.Api, error) {
	s := battle.NewService()

	api := grpc.NewApi()

	api.RegisterService(&battlebotspb.Battle_ServiceDesc, s)

	return api, nil
}
