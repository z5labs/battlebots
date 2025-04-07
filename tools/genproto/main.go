package main

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
)

func main() {
	sigCtx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer stop()

	os.Exit(run(sigCtx))
}

func run(ctx context.Context) (code int) {
	log := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}))

	args := os.Args
	if len(args) < 2 {
		log.ErrorContext(ctx, "must provide folder")
		return 1
	}

	folderPath := args[1]
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		log.ErrorContext(ctx, "failed to list files in directory", slog.String("error", err.Error()))
		return 1
	}

	protocArgs := []string{
		"-I", folderPath,
		"--go_out", ".",
		"--go_opt", "paths=source_relative",
		"--go_opt", "default_api_level=API_OPAQUE",
		"--go-grpc_out", ".",
		"--go-grpc_opt", "paths=source_relative",
	}
	for _, entry := range entries {
		protocArgs = append(protocArgs, filepath.Join(folderPath, entry.Name()))
	}

	cmd := exec.CommandContext(ctx, "protoc", protocArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.ErrorContext(ctx, "failed to run protoc", slog.String("error", err.Error()))
		return 1
	}
	return
}
