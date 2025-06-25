package simple_algo

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	signalModel "github.com/kirill-a-belov/trader/internal/signal/model"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

func (m *Module) Process(ctx context.Context) error {
	_, span := tracer.Start(ctx, "pkg.internal.simple_algo.Process")
	defer span.End()

	if err := m.preload(ctx); err != nil {
		return errors.Wrap(err, "preload")
	}

	priceFeed, err := m.bybitModule.PriceFeedBTCUSDT(ctx)
	if err != nil {
		return errors.Wrap(err, "bybitModule.PriceFeedBTCUSDT")
	}

	go func() {
		for priceTicker := range priceFeed {
			sl := signalModel.SignalList{}
			for _, candleModule := range m.candleModuleList {
				signalList, err := candleModule.Insert(ctx, priceTicker)
				if err != nil {
					m.log.Error(ctx, "inserting candle: ",
						"error", err,
						"priceTicker", priceTicker,
						"timeframe", candleModule.Config.Timeframe,
					)
				}
				sl = append(sl, signalList...)
			}

			fmt.Printf("\rLast ticker (ts - lp): %v - %v", priceTicker.Timestamp, priceTicker.LastPrice)

			if priceTicker.IsTooOld() {
				continue
			}

			if err := m.managePosition(ctx, priceTicker, sl); err != nil {
				m.log.Error(ctx, "managePosition",
					"error", err,
				)

				continue
			}
		}
	}()

	return nil
}
