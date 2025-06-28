package main

import (
	"bytes"
	_ "embed"

	"github.com/z5labs/battlebots/services/battle/app"

	"github.com/z5labs/humus/grpc"
)

//go:embed config.yaml
var configBytes []byte

func main() {
	grpc.Run(bytes.NewReader(configBytes), app.Init)
}
