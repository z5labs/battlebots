package app

import (
	"context"
	crand "crypto/rand"
	"io"
	"time"

	"github.com/z5labs/battlebots/cmd/poc/bot"
	"github.com/z5labs/battlebots/pkgs/battlebotspb"

	"github.com/z5labs/humus/job"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	job.Config `config:",squash"`

	BattleService struct {
		Endpoint string `config:"endpoint"`
	} `config:"battle_service"`

	Bot struct {
		WaitFor struct {
			Min time.Duration `config:"min"`
			Max time.Duration `config:"max"`
		} `config:"wait_for"`
	} `config:"bot"`
}

func Init(ctx context.Context, cfg Config) (*job.App, error) {
	conn, err := grpc.NewClient(
		cfg.BattleService.Endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	var seed [32]byte
	_, err = io.ReadFull(crand.Reader, seed[:])
	if err != nil {
		return nil, err
	}

	h, err := bot.NewHandler(
		battlebotspb.NewBattleClient(conn),
		bot.MinWaitFor(cfg.Bot.WaitFor.Min),
		bot.MaxWaitFor(cfg.Bot.WaitFor.Max),
	)
	if err != nil {
		return nil, err
	}

	app := job.NewApp(h)
	return app, nil
}
