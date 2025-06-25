package candle

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/kirill-a-belov/trader/internal/candle/model"
	chartModel "github.com/kirill-a-belov/trader/internal/chart/model"
	"github.com/kirill-a-belov/trader/pkg/bybit"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

func (m *Module) Preload(ctx context.Context, candleList []*bybit.Kline) error {
	_, span := tracer.Start(ctx, "pkg.internal.candle.Preload")
	defer span.End()

	timeframe := m.Config.Timeframe
	now := time.Now()
	parts := int(bybit.KlineTimeframe / timeframe)

	for i := len(candleList) - 1; i >= 0; i-- {
		candle := candleList[i]
		if candle.StartTime.Before(now.Add(-m.Config.Depth)) {
			continue
		}

		tsList := []time.Time{}
		switch {
		case timeframe >= bybit.KlineTimeframe:
			tsList = append(tsList, candle.StartTime.Truncate(timeframe))
		case timeframe < bybit.KlineTimeframe:
			for i := 0; i < parts; i++ {
				tsList = append(tsList, candle.StartTime.Add(time.Duration(i)*timeframe).Truncate(timeframe))
			}
		}

		for _, ts := range tsList {
			if _, ok := m.candleStorage[ts]; !ok {
				c := &model.Candle{
					Timestamp: ts,
					High:      candle.High,
					Low:       candle.Low,
					Open:      candle.Open,
					Close:     candle.Close,
				}
				m.candleStorage[ts] = c

				if !m.Config.SendCandleToChart {
					continue
				}

				if err := m.chartModule.PutCandle(ctx, &chartModel.Candle{
					Time:  c.Timestamp.Unix(),
					Open:  c.Open,
					High:  c.High,
					Low:   c.Low,
					Close: c.Close,
				}); err != nil {
					return errors.Wrap(err, "put candle to chart")
				}
			}
		}
	}

	return nil
}
