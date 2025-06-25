package signal

import (
	"context"

	"github.com/pkg/errors"

	candleModel "github.com/kirill-a-belov/trader/internal/candle/model"
	"github.com/kirill-a-belov/trader/internal/signal/model"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

func (m *Module) SignalList(
	ctx context.Context,
	candles []*candleModel.Candle,
) (model.SignalList, error) {
	_, span := tracer.Start(ctx, "pkg.internal.signal.SignalList")
	defer span.End()

	signals := []*model.Signal{}

	if candles == nil || len(candles) < m.CandleListNecessaryDepth() {
		signals = append(signals, m.signal(ctx, "not enough candles", model.SignalDirectionFlat, 0))

		return signals, nil
	}

	emaSignalList, err := m.emaSignalGroup(ctx, candles)
	if err != nil {
		return nil, errors.Wrap(err, "ema signal list")
	}
	signals = append(signals, emaSignalList...)

	priceSignalList, err := m.priceSignalGroup(ctx, candles)
	if err != nil {
		return nil, errors.Wrap(err, "ema signal list")
	}
	signals = append(signals, priceSignalList...)

	if len(signals) == 0 {
		signals = append(signals, m.signal(ctx, "no signals", model.SignalDirectionFlat, 0))
	}

	return signals, nil
}
