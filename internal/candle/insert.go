package candle

import (
	"context"

	"github.com/pkg/errors"

	"github.com/kirill-a-belov/trader/internal/candle/model"
	chartModel "github.com/kirill-a-belov/trader/internal/chart/model"
	signalModel "github.com/kirill-a-belov/trader/internal/signal/model"
	"github.com/kirill-a-belov/trader/pkg/ticker"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

func (m *Module) sendCandleToChart(ctx context.Context, candle *model.Candle, signalList signalModel.SignalList) error {
	_, span := tracer.Start(ctx, "pkg.internal.candle.sendCandleToChart")
	defer span.End()

	if m.Config.SendCandleToChart {
		if err := m.chartModule.PutCandle(ctx, &chartModel.Candle{
			Time:       candle.Timestamp.Unix(),
			Open:       candle.Open,
			High:       candle.High,
			Low:        candle.Low,
			Close:      candle.Close,
			BuyVolume:  candle.BuyVolume,
			SellVolume: candle.SellVolume,
			BuyCount:   candle.BuyCount,
			SellCount:  candle.SellCount,
			TickCount:  candle.TickCount,
			TickSum:    candle.TickSum,
			TradeList:  candle.TradeList,
		}); err != nil {
			return errors.Wrap(err, "put candle to chart")
		}
	}

	if m.Config.SendMarkerToChart {
		for _, signal := range signalList {
			if err := m.chartModule.PutMarker(ctx, m.chartModule.FromSignal(signal)); err != nil {
				m.log.Error(ctx, "failed to put marker in chart",
					"error", err,
					"signal", signal,
				)

				return errors.New("failed to put marker in chart")
			}
		}
	}

	return nil
}

func (m *Module) Insert(ctx context.Context, priceTicker *ticker.Ticker) (signalModel.SignalList, error) {
	_, span := tracer.Start(ctx, "pkg.internal.candle.Insert")
	defer span.End()

	candleTime := priceTicker.Timestamp.Truncate(m.Config.Timeframe)
	delete(m.candleStorage, candleTime.Add(-m.Config.Depth))

	candle, ok := m.candleStorage[candleTime]
	if !ok {
		candle = &model.Candle{
			Timestamp: candleTime,
			Open:      priceTicker.LastPrice,
			High:      priceTicker.LastPrice,
			Low:       priceTicker.LastPrice,
		}
		m.candleStorage[candleTime] = candle
	}
	{
		if priceTicker.LastPrice > candle.High {
			candle.High = priceTicker.LastPrice
		}
		if priceTicker.LastPrice < candle.Low {
			candle.Low = priceTicker.LastPrice
		}
		candle.Close = priceTicker.LastPrice
		candle.TickCount++
		candle.TickSum += priceTicker.LastPrice

		candle.BuyCount += priceTicker.BuyCount
		candle.SellCount += priceTicker.SellCount
		candle.BuyVolume += priceTicker.BuyVolume
		candle.SellVolume += priceTicker.SellVolume
		candle.TradeList = append(candle.TradeList, priceTicker.TradeList...)
	}

	signalList, err := m.signalModule.SignalList(ctx, m.candleListOldToNew(ctx, candle.Timestamp, m.signalModule.CandleListNecessaryDepth()))
	if err != nil {
		return nil, errors.Wrap(err, "insert signal list")
	}

	if err := m.sendCandleToChart(ctx, candle, signalList); err != nil {
		m.log.Error(ctx, "failed to send candle to chart",
			"error", err,
			"candle", candle,
		)
	}

	return signalList, nil
}
