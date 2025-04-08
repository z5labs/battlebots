package cmd

import (
	"context"
	crand "crypto/rand"
	"io"
	"log/slog"
	"math/rand/v2"
	"os"
	"os/signal"
	"time"

	"github.com/z5labs/battlebots/pkgs/battlebotspb"
	"github.com/z5labs/battlebots/pkgs/poc"

	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	experimental "google.golang.org/grpc/experimental/opentelemetry"
	"google.golang.org/grpc/stats/opentelemetry"
)

func ExitWithCodeOnError(code int, err error) {
	if err == nil {
		return
	}
	os.Exit(code)
}

func Run(ctx context.Context) error {
	sigCtx, cancel := signal.NotifyContext(ctx, os.Kill, os.Interrupt)
	defer cancel()

	cmd := &cobra.Command{
		Use:     "poc",
		PreRun:  initOTelSDK,
		Run:     run,
		PostRun: shutdownOTelSDK,
	}

	cmd.Flags().String("otlp-endpoint", "", "Specify OTLP endpoint for sending telemetry signals to.")
	cmd.Flags().String("game-server-endpoint", "", "Specify the game server endpoint.")
	cmd.Flags().Duration("min-wait-for", 16*time.Millisecond, "Set the minimum time to wait before randomly moving.")
	cmd.Flags().Duration("max-wait-for", 100*time.Millisecond, "Set the maximum time to wait before randomly moving.")

	cmd.MarkFlagRequired("otlp-endpoint")
	cmd.MarkFlagRequired("game-server-endpoint")

	return cmd.ExecuteContext(sigCtx)
}

func initOTelSDK(cmd *cobra.Command, args []string) {
	// TODO
}

func shutdownOTelSDK(cmd *cobra.Command, args []string) {
	// TODO
}

func run(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	log := otelslog.NewLogger("poc")

	minWaitFor, err := cmd.Flags().GetDuration("min-wait-for")
	if err != nil {
		log.ErrorContext(ctx, "failed to get min wait for duration", slog.String("error", err.Error()))
		return
	}

	maxWaitFor, err := cmd.Flags().GetDuration("max-wait-for")
	if err != nil {
		log.ErrorContext(ctx, "failed to get max wait for duration", slog.String("error", err.Error()))
		return
	}

	gameServerEndpoint, err := cmd.Flags().GetString("game-server-endpoint")
	if err != nil {
		log.ErrorContext(ctx, "failed to get game server endpoint", slog.String("error", err.Error()))
		return
	}

	cc, err := grpc.NewClient(
		gameServerEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		opentelemetry.DialOption(opentelemetry.Options{
			MetricsOptions: opentelemetry.MetricsOptions{
				MeterProvider: otel.GetMeterProvider(),
			},
			TraceOptions: experimental.TraceOptions{
				TracerProvider:    otel.GetTracerProvider(),
				TextMapPropagator: otel.GetTextMapPropagator(),
			},
		}),
	)
	if err != nil {
		log.ErrorContext(ctx, "failed to init grpc conn", slog.String("error", err.Error()))
		return
	}

	battleClient := battlebotspb.NewBattleClient(cc)

	var seed [32]byte
	_, err = io.ReadFull(crand.Reader, seed[:])
	if err != nil {
		log.ErrorContext(ctx, "failed to read crypto bytes", slog.String("error", err.Error()))
		return
	}

	bot := poc.NewBot(
		battleClient,
		rand.NewChaCha8(seed),
		poc.MinWaitFor(minWaitFor),
		poc.MaxWaitFor(maxWaitFor),
	)

	err = bot.Run(ctx)
	if err != nil {
		log.ErrorContext(ctx, "failed to run bot", slog.String("error", err.Error()))
		return
	}
}
