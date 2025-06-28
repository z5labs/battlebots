package main

import (
	"bytes"
	_ "embed"

	"github.com/z5labs/battlebots/example/random-move-bot-go/bot"
	"github.com/z5labs/battlebots/sdk/battlebots-go"
)

//go:embed config.yaml
var configBytes []byte

func main() {
	battlebots.Run(bytes.NewReader(configBytes), bot.Init)
}
