package trader

import (
	"context"

	"github.com/pkg/errors"

	"github.com/kirill-a-belov/trader/pkg/tracer"
)

func (m *Module) Process(ctx context.Context) error {
	_, span := tracer.Start(ctx, "pkg.internal.trader.Process")
	defer span.End()

	if err := m.chartModule.Serve(ctx); err != nil {
		return errors.Wrap(err, "chartModule.Serve")
	}

	for _, alg := range m.tradingAlgorithmList {
		if err := alg.Process(context.Background()); err != nil {
			return errors.Wrapf(err, "failed to process trader module (%s)", alg.Name())
		}

		m.log.Info(ctx, "processed trader module",
			"name", alg.Name(),
		)
	}

	return nil
}
