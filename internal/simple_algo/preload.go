package simple_algo

import (
	"context"

	"github.com/pkg/errors"

	"github.com/kirill-a-belov/trader/pkg/tracer"
)

func (m *Module) preload(ctx context.Context) error {
	_, span := tracer.Start(ctx, "pkg.internal.simple_algo.preload")
	defer span.End()

	candleList, err := m.bybitModule.CandleListBTCUSDT(ctx)
	if err != nil {
		return errors.Wrap(err, "bybitModule.CandleListBTCUSDT")
	}

	for _, candleModule := range m.candleModuleList {
		if err := candleModule.Preload(ctx, candleList); err != nil {
			return errors.Wrap(err, "candleModule.Preload")
		}
	}

	return nil
}
