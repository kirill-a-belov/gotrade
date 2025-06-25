package position

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	chartModel "github.com/kirill-a-belov/trader/internal/chart/model"
	"github.com/kirill-a-belov/trader/internal/position/model"
	"github.com/kirill-a-belov/trader/pkg/bybit"
	"github.com/kirill-a-belov/trader/pkg/ptr"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

func (m *Module) OpenPosition(ctx context.Context, price float64, positionSide model.PositionSide, comment string) error {
	_, span := tracer.Start(ctx, "pkg.internal.position.OpenPosition")
	defer span.End()

	if m.currentPosition != nil {
		return errors.New("position already opened")
	}

	positionStopLossPrice := m.positionStopLossPrice(price, positionSide)

	if m.config.Debug {
		fmt.Println()
		fmt.Println("--------Open-Position---------:", positionSide)
		fmt.Println("Calculation details:")
		fmt.Println("price", price)
		fmt.Println("stopLossDiff:", price-positionStopLossPrice)
		fmt.Println("positionStopLossPrice ", positionStopLossPrice)
		fmt.Println("comment", comment)
		fmt.Println("------------------------------")
		fmt.Println()
	}

	orderRequest := &bybit.CreateOrderRequest{
		Category:    "linear",
		Symbol:      "BTCUSDT",
		Side:        string(positionSide),
		OrderType:   "Market",
		Qty:         fmt.Sprintf("%f", m.config.PositionSize),
		TimeInForce: "ImmediateOrCancel",
		StopLoss:    ptr.PtrString(fmt.Sprintf("%f", positionStopLossPrice)),
		SlTriggerBy: ptr.PtrString("LastPrice"),
	}

	_, err := m.bybitModule.CreateOrder(ctx, orderRequest)
	if err != nil {
		return errors.Wrap(err, "bybitModule.CreateOrder")
	}

	now := time.Now()

	m.currentPosition = &model.Position{
		Id:            uuid.NewString(),
		CreatedAt:     now,
		Side:          positionSide,
		Size:          m.config.PositionSize,
		StopLossPrice: positionStopLossPrice,
	}

	var colour string
	switch positionSide {
	case model.PositionSideBuy:
		colour = "green"
	case model.PositionSideSell:
		colour = "red"
	}

	if err := m.chartModule.PutMarker(ctx, &chartModel.Marker{
		Name:      "open position",
		Timestamp: now,
		Timeframe: time.Minute,
		Time:      now.Unix(),
		Text:      "✚",
		Color:     colour,
		Position:  "belowBar",
		Details: fmt.Sprintf("%s | %s | %s | %s",
			now.Format("15:04:05"), "POSITION: open", m.currentPosition.Side, comment),
	}); err != nil {
		m.log.Error(ctx, "chartModule.PutMarker",
			"err", err,
		)
	}

	return nil
}
