package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/kirill-a-belov/trader/pkg/tracer"
)

func New(ctx context.Context) *cobra.Command {
	span, _ := tracer.Start(ctx, "cmd.New")
	defer span.Done()

	cmd := &cobra.Command{
		Short: "Trading platform",
	}
	cmd.AddCommand(
		startCMD(ctx),
	)

	return cmd
}
