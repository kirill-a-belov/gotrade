package signal

import (
	"context"
	candleModel "github.com/kirill-a-belov/trader/internal/candle/model"
	"github.com/kirill-a-belov/trader/internal/signal/model"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

const (
	SignalNameHugeBullCandle = "huge bull candle"
	SignalNameHugeBearCandle = "huge bear candle"
)

func (m *Module) priceSignalGroup(ctx context.Context, candles []*candleModel.Candle) (model.SignalList, error) {
	_, span := tracer.Start(ctx, "pkg.internal.signal.priceSignalGroup")
	defer span.End()

	lastCandleIndex := len(candles) - 1

	signalList := []*model.Signal{}

	lastCandleHugeBull := candles[lastCandleIndex].Close-candles[lastCandleIndex].Open >
		candles[lastCandleIndex].High*float64(m.Config.PriceHugeCandlePercent)/100
	secondLastCandleHugeBull := candles[lastCandleIndex-1].Close-candles[lastCandleIndex-1].Open >
		candles[lastCandleIndex-1].High*float64(m.Config.PriceHugeCandlePercent)/100
	if lastCandleHugeBull && !secondLastCandleHugeBull {
		signalList = append(signalList, m.signal(ctx, SignalNameHugeBullCandle, model.SignalDirectionBear, 1))
	}

	lastCandleHugeBear := candles[lastCandleIndex].Open-candles[lastCandleIndex].Close >
		candles[lastCandleIndex].High*float64(m.Config.PriceHugeCandlePercent)/100
	secondLastCandleHugeBear := candles[lastCandleIndex-1].Open-candles[lastCandleIndex-1].Close >
		candles[lastCandleIndex-1].High*float64(m.Config.PriceHugeCandlePercent)/100
	if lastCandleHugeBear && !secondLastCandleHugeBear {
		signalList = append(signalList, m.signal(ctx, SignalNameHugeBearCandle, model.SignalDirectionBull, 1))
	}

	return signalList, nil
}
