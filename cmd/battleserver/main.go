// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/z5labs/battlebots/internal/battle2d"
	"github.com/z5labs/battlebots/internal/battle3d"
	"github.com/z5labs/battlebots/pkgs/battlepb"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	s := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	battlepb.RegisterBattle2DServer(s, battle2d.NewServer())
	battlepb.RegisterBattle3DServer(s, battle3d.NewServer())

	ls, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.ErrorContext(
			ctx,
			"failed to listen for connections",
			slog.Int("port", 9090),
			slog.String("error", err.Error()),
		)
		return
	}

	eg, egctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return s.Serve(ls)
	})
	eg.Go(func() error {
		<-egctx.Done()
		s.GracefulStop()
		return nil
	})

	err = eg.Wait()
	if err == nil {
		return
	}
	log.Error("unexpected error while serving", slog.String("error", err.Error()))
}
