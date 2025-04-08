package cmd

import (
	"context"
	crand "crypto/rand"
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"os"
	"os/signal"
	"time"

	"github.com/z5labs/battlebots/pkgs/battlebotspb"

	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		Run:     poc,
		PostRun: shutdownOTelSDK,
	}

	cmd.Flags().String("otlp-endpoint", "", "Specify OTLP endpoint for sending telemetry signals to.")
	cmd.Flags().String("game-server-endpoint", "", "Specify the game server endpoint.")

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

func poc(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	log := otelslog.NewLogger("poc")

	// TODO: connect to game server
	gameServerEndpoint, err := cmd.Flags().GetString("game-server-endpoint")
	if err != nil {
		log.ErrorContext(ctx, "failed to get game server endpoint", slog.String("error", err.Error()))
		return
	}

	cc, err := grpc.NewClient(gameServerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.ErrorContext(ctx, "failed to init grpc conn", slog.String("error", err.Error()))
		return
	}

	battlebots := battlebotspb.NewBattleClient(cc)

	eg, egctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		var subscription battlebotspb.StateChangeSubscription
		events, err := battlebots.State(egctx, &subscription)
		if err != nil {
			return err
		}

		for {
			event, err := events.Recv()
			if err != nil {
				return err
			}

			processEvent(egctx, log, event)
		}
	})

	eg.Go(func() error {
		var seed [32]byte
		_, err := io.ReadFull(crand.Reader, seed[:])
		if err != nil {
			return err
		}

		r := rand.New(rand.NewChaCha8(seed))

		fmt.Println("Press the enter key to begin randomly moving...")
		_, err = fmt.Scanln()
		if err != nil {
			return err
		}

		for {
			// 60 FPS = 16.6666... ms per frame so we test with double frame time
			// as worst case for delay between move actions. This was selected completed
			// arbitrarily.
			waitFor := time.Duration(r.IntN(32)) * time.Millisecond
			select {
			case <-egctx.Done():
				return egctx.Err()
			case <-time.After(waitFor):
			}

			err = randomlyMove(egctx, battlebots)
			if err != nil {
				return err
			}
		}
	})
	err = eg.Wait()
	if err != nil {
		log.ErrorContext(ctx, "unexpected error while running the poc", slog.String("error", err.Error()))
		return
	}
}

func processEvent(ctx context.Context, log *slog.Logger, event *battlebotspb.StateChangeEvent) {
	spanCtx, span := otel.Tracer("cmd").Start(ctx, "processEvent")
	defer span.End()

	if !event.HasPosition() {
		log.WarnContext(spanCtx, "expected to only receive position state changes")
		return
	}

	bot := event.GetBot()
	pos := event.GetPosition()

	log.InfoContext(
		spanCtx,
		"received position state change event",
		slog.String("bot.id", bot.GetId()),
		slog.Float64("position.x0", pos.GetX0()),
		slog.Float64("position.x1", pos.GetX1()),
	)
}

func randomlyMove(ctx context.Context, battlebots battlebotspb.BattleClient) error {
	spanCtx, span := otel.Tracer("cmd").Start(ctx, "randomlyMove")
	defer span.End()

	var velocity battlebotspb.Vector
	velocity.SetX0(rand.Float64())
	velocity.SetX1(rand.Float64())

	var req battlebotspb.MoveRequest
	req.SetVelocity(&velocity)

	_, err := battlebots.Move(spanCtx, &req)
	return err
}
