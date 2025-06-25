package signal

import (
	"context"
	"fmt"
	candleModel "github.com/kirill-a-belov/trader/internal/candle/model"
	"github.com/kirill-a-belov/trader/internal/signal/model"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

const (
	SignalNameTrendGoesUp      = "trend goes up"
	SignalNameTrendGoesDown    = "trend goes down"
	SignalNameBullEmaCrossover = "fast and slow EMA crossover"
	SignalNameBearEmaCrossover = "fast and slow EMA crossover"
)

func (m *Module) emaSignalGroup(ctx context.Context, candles []*candleModel.Candle) (model.SignalList, error) {
	_, span := tracer.Start(ctx, "pkg.internal.signal.emaSignalGroup")
	defer span.End()

	fastEMA := ema(ctx, candles, m.Config.FastEmaPeriod)
	slowEMA := ema(ctx, candles, m.Config.SlowEmaPeriod)
	if len(fastEMA) < len(candles) || len(slowEMA) < len(candles) {
		return nil, fmt.Errorf("EMA calculation failed")
	}

	lastCandleIndex := len(candles) - 1
	prevCandleIndex := lastCandleIndex - 1

	signalList := []*model.Signal{}
	if fastEMA[prevCandleIndex] < slowEMA[prevCandleIndex] && fastEMA[lastCandleIndex] > slowEMA[lastCandleIndex] {
		signalList = append(signalList, m.signal(ctx, SignalNameBullEmaCrossover, model.SignalDirectionBull, 1))
	}
	if fastEMA[prevCandleIndex] > slowEMA[prevCandleIndex] && fastEMA[lastCandleIndex] < slowEMA[lastCandleIndex] {
		signalList = append(signalList, m.signal(ctx, SignalNameBearEmaCrossover, model.SignalDirectionBear, 1))
	}
	if fastEMA[lastCandleIndex] > slowEMA[lastCandleIndex] {
		signalList = append(signalList, m.signal(ctx, SignalNameTrendGoesUp, model.SignalDirectionBull, 1))
	}
	if fastEMA[lastCandleIndex] < slowEMA[lastCandleIndex] {
		signalList = append(signalList, m.signal(ctx, SignalNameTrendGoesDown, model.SignalDirectionBear, 1))
	}

	return signalList, nil
}

func ema(ctx context.Context, candles []*candleModel.Candle, period int) []float64 {
	_, span := tracer.Start(ctx, "pkg.internal.signal.ema")
	defer span.End()

	if len(candles) < period {
		return nil
	}

	result := make([]float64, len(candles))
	alpha := 2.0 / float64(period+1)

	sum := 0.0
	for i := 0; i < period; i++ {
		sum += candles[i].Close
	}
	result[period-1] = sum / float64(period)

	for i := period; i < len(candles); i++ {
		result[i] = alpha*candles[i].Close + (1-alpha)*result[i-1]
	}

	return result
}
