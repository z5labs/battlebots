package main

import (
	"context"

	"github.com/z5labs/battlebots/pkgs/poc/cmd"
)

func main() {
	cmd.ExitWithCodeOnError(1, cmd.Run(context.Background()))
}
