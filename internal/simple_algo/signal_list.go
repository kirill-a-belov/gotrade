package simple_algo

import (
	"context"
	"time"

	positionModel "github.com/kirill-a-belov/trader/internal/position/model"
	signal "github.com/kirill-a-belov/trader/internal/signal"
	signalModel "github.com/kirill-a-belov/trader/internal/signal/model"
	"github.com/kirill-a-belov/trader/pkg/ticker"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

func (m *Module) canOpenPositionBySignalList(ctx context.Context, signalList signalModel.SignalList) (positionModel.PositionSide, string, error) {
	_, span := tracer.Start(ctx, "pkg.internal.simple_algo.canOpenPositionBySignalList")
	defer span.End()

	signalMap := signalList.TimeframeAndNameMap()
	_, hugeBullCandle1m := signalMap[signalModel.TimeframeAndName(time.Minute, signal.SignalNameHugeBullCandle)]
	_, hugeBearCandle1m := signalMap[signalModel.TimeframeAndName(time.Minute, signal.SignalNameHugeBearCandle)]
	_, bullTrend1m := signalMap[signalModel.TimeframeAndName(time.Minute, signal.SignalNameTrendGoesUp)]
	_, bullTrend10s := signalMap[signalModel.TimeframeAndName(10*time.Second, signal.SignalNameTrendGoesUp)]
	_, bullTrend15s := signalMap[signalModel.TimeframeAndName(15*time.Second, signal.SignalNameTrendGoesUp)]
	_, bearTrend1m := signalMap[signalModel.TimeframeAndName(time.Minute, signal.SignalNameTrendGoesDown)]
	_, bearTrend10s := signalMap[signalModel.TimeframeAndName(10*time.Second, signal.SignalNameTrendGoesDown)]
	_, bearTrend15s := signalMap[signalModel.TimeframeAndName(15*time.Second, signal.SignalNameTrendGoesDown)]

	if hugeBullCandle1m && !bullTrend1m {
		return positionModel.PositionSideSell, "sell after huge bull candle without trend", nil
	}
	if hugeBearCandle1m && !bearTrend1m {
		return positionModel.PositionSideBuy, "sell after huge bear candle without trend", nil
	}

	if bullTrend1m && bullTrend10s && bullTrend15s {
		// TODO flat market - case
		return positionModel.PositionSideBuy, "buy in the beginning of bull trend", nil
	}

	if bearTrend10s && bearTrend15s && bearTrend1m {
		return positionModel.PositionSideSell, "sell in the beginning of bear trend", nil
	}

	// TODO More cases

	return positionModel.PositionSideUnknown, "", nil
}

func (m *Module) shouldClosePositionBySignalLists(
	ctx context.Context,
	signalList signalModel.SignalList,
	priceTicker *ticker.Ticker,
	position *positionModel.Position,
) (
	bool,
	string, error,
) {
	_, span := tracer.Start(ctx, "pkg.internal.simple_algo.shouldClosePositionBySignalLists")
	defer span.End()

	signalMap := signalList.TimeframeAndNameMap()
	_, bullTrend1m := signalMap[signalModel.TimeframeAndName(time.Minute, signal.SignalNameTrendGoesUp)]
	_, bullTrend10s := signalMap[signalModel.TimeframeAndName(10*time.Second, signal.SignalNameTrendGoesUp)]
	_, bearTrend1m := signalMap[signalModel.TimeframeAndName(time.Minute, signal.SignalNameTrendGoesDown)]
	_, bearTrend10s := signalMap[signalModel.TimeframeAndName(10*time.Second, signal.SignalNameTrendGoesDown)]

	clearPnl := m.clearPNL(ctx, priceTicker, position)

	if position.Side == positionModel.PositionSideSell &&
		bearTrend1m &&
		!bearTrend10s && clearPnl > 0.5 {

		return true, "closing by bear trend reversing and small positive profit", nil
	}
	if position.Side == positionModel.PositionSideBuy &&
		bullTrend1m &&
		!bullTrend10s && clearPnl > 0.5 {

		return true, "closing by bull trend reversing and small positive profit", nil
	}

	// TODO More cases

	return false, "", nil
}
