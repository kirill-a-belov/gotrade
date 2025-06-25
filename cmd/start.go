package cmd

import (
	"context"
	"github.com/kirill-a-belov/trader/internal/trader"
	"github.com/kirill-a-belov/trader/pkg/logger"
	"github.com/pkg/errors"

	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/kirill-a-belov/trader/pkg/tracer"
)

func startCMD(ctx context.Context) *cobra.Command {
	ctx, span := tracer.Start(ctx, "cmd.startCMD")
	defer span.End()

	return &cobra.Command{
		Use:   "start",
		Short: "Start trading algorithms",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, span := tracer.Start(ctx, "cmd.api.New")
			defer span.End()

			traderModule, err := trader.New(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to create trader module")
			}

			if err := traderModule.Process(ctx); err != nil {
				return errors.Wrap(err, "failed to process trader module")
			}

			log := logger.New("api_server_cmd")
			log.Info(ctx, "Trader platform module started")

			waitSigterm()
			softStop(ctx, traderModule, log)

			return nil
		},
	}
}

func waitSigterm() {
	signals := make(chan os.Signal, 8)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	<-signals
}

func softStop(ctx context.Context, platform interface{}, log logger.Logger) {
	if platform == nil {
		os.Exit(0)
	}

	log.Info(ctx, "waiting for platform shutdown")

	//TODO Sell all positions
	/*
		ctxTime, _ := context.WithTimeout(context.Background(), time.Second*5)
		if err := server.Shutdown(ctxTime); err != nil {
			log.Error(ctx, "failed to shutdown platform, forcing to close", "err", err)

			if err = server.Close(); err != nil {
				log.Error(ctx, "close failed", "err", err)
			}
		}*/

	os.Exit(0)
}
