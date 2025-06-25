package main

import (
	"context"

	"github.com/kirill-a-belov/trader/cmd"
)

func main() {
	ctx := context.Background()

	cli := cmd.New(ctx)
	if err := cli.Execute(); err != nil {
		panic(err)
	}
}
