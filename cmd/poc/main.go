package main

import (
	"bytes"
	_ "embed"

	"github.com/z5labs/battlebots/cmd/poc/app"

	"github.com/z5labs/humus/job"
)

//go:embed config.yaml
var configBytes []byte

func main() {
	job.Run(bytes.NewReader(configBytes), app.Init)
}
