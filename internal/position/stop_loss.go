package position

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	chartModel "github.com/kirill-a-belov/trader/internal/chart/model"
	"github.com/kirill-a-belov/trader/internal/position/model"
	"github.com/kirill-a-belov/trader/pkg/bybit"
	"github.com/kirill-a-belov/trader/pkg/ptr"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

func (m Module) positionStopLossPrice(price float64, side model.PositionSide) float64 {
	var (
		positionStopLossPrice float64
		stopLossDiff          float64
	)
	{
		stopLossDiff = m.config.MaxRiskSizeUSDT / m.config.PositionSize
		switch side {
		case model.PositionSideSell:
			positionStopLossPrice = price + stopLossDiff
		case model.PositionSideBuy:
			positionStopLossPrice = price - stopLossDiff
		}
	}

	return positionStopLossPrice
}

func (m *Module) UpdatePositionStopLoss(ctx context.Context, price float64) error {
	_, span := tracer.Start(ctx, "pkg.internal.position.UpdatePositionStopLoss")
	defer span.End()

	if m.currentPosition == nil {
		return errors.New("position not opened")
	}

	oldStopLossPrice := m.currentPosition.StopLossPrice
	newStopLossPrice := m.positionStopLossPrice(price, m.currentPosition.Side)

	if (m.currentPosition.Side == model.PositionSideBuy && newStopLossPrice <= oldStopLossPrice) ||
		(m.currentPosition.Side == model.PositionSideSell && newStopLossPrice >= oldStopLossPrice) {
		return nil
	}

	if m.config.Debug {
		fmt.Println()
		fmt.Println("--------UpdateStopLoss-Position---------:")
		fmt.Println("Calculation details:")
		fmt.Println("price", price)
		fmt.Println("oldStopLossPrice:", oldStopLossPrice)
		fmt.Println("newStopLossPrice ", newStopLossPrice)
		fmt.Println("------------------------------")
		fmt.Println()
	}

	orderRequest := &bybit.SetTradeConditionsBTCUSDTRequest{
		Category:    "linear",
		Symbol:      "BTCUSDT",
		TpslMode:    "Full",
		PositionIdx: 0,
		StopLoss:    ptr.PtrString(fmt.Sprintf("%f", newStopLossPrice)),
		SlTriggerBy: ptr.PtrString("LastPrice"),
	}

	err := m.bybitModule.SetTradeConditionsBTCUSDT(ctx, orderRequest)
	if err != nil {
		return errors.Wrap(err, "bybitModule.SetTradeConditionsBTCUSDT")
	}

	m.currentPosition.StopLossPrice = newStopLossPrice

	now := time.Now()
	if err := m.chartModule.PutMarker(ctx, &chartModel.Marker{
		Name:      "position stop loss adjustment",
		Timestamp: now,
		Time:      now.Unix(),
		Timeframe: time.Minute,
		Text:      "●",
		Color:     "grey",
		Position:  "belowBar",
		Details: fmt.Sprintf("%s | %s | %s | %s",
			now.Format("15:04:05"), "POSITION: adjust", m.currentPosition.Side, "change stop loss price"),
	}); err != nil {
		m.log.Error(ctx, "chartModule.PutMarker",
			"err", err,
		)
	}

	return nil
}
