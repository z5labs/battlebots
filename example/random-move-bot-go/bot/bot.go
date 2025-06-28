package bot

import (
	"context"
	"log/slog"

	"github.com/z5labs/battlebots/sdk/battlebots-go"
	"github.com/z5labs/humus"
)

type Config struct {
	battlebots.Config `config:",squash"`
}

type Controller struct {
	log *slog.Logger
	bot *battlebots.Bot
}

func Init(ctx context.Context, cfg Config, bot *battlebots.Bot) (*Controller, error) {
	c := &Controller{
		log: humus.Logger("bot"),
		bot: bot,
	}

	return c, nil
}

func (c *Controller) Tick(ctx context.Context, state *battlebots.DefaultState) error {
	c.log.InfoContext(ctx, "tick")
	return nil
}
