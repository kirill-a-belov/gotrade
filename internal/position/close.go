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

func (m *Module) ClosePosition(ctx context.Context, comment string) error {
	_, span := tracer.Start(ctx, "pkg.internal.position.ClosePosition")
	defer span.End()

	if m.currentPosition == nil {
		return errors.New("position not opened")
	}

	if m.config.Debug {
		fmt.Println()
		fmt.Println("--------Close-Position---------:", m.currentPosition.Side)
		fmt.Println("comment: ", comment)
		fmt.Println()
	}

	orderRequest := &bybit.CreateOrderRequest{
		Category:    "linear",
		Symbol:      "BTCUSDT",
		OrderType:   "Market",
		Qty:         fmt.Sprintf("%f", m.currentPosition.Size),
		TimeInForce: "ImmediateOrCancel",
		ReduceOnly:  ptr.PtrBool(false),
		Side:        string(m.currentPosition.OppositeSide()),
	}

	const (
		closeAttemptPause   = 500 * time.Millisecond
		closeAttemptMaxSize = 10
	)
	var closeAttemptCounter int
	for {
		if closeAttemptCounter >= closeAttemptMaxSize {
			panic(errors.New("too many attempts to close position"))
		}

		_, err := m.bybitModule.CreateOrder(ctx, orderRequest)
		if err != nil {
			m.log.Error(ctx, "error closing position request", err)
			closeAttemptCounter++
			time.Sleep(closeAttemptPause)

			continue
		}

		break
	}
	now := time.Now()

	var colour string
	switch m.currentPosition.Side {
	case model.PositionSideBuy:
		colour = "green"
	case model.PositionSideSell:
		colour = "red"
	}
	if err := m.chartModule.PutMarker(ctx, &chartModel.Marker{
		Name:      "position close",
		Timestamp: now,
		Time:      now.Unix(),
		Timeframe: time.Minute,
		Text:      "━",
		Color:     colour,
		Position:  "belowBar",
		Details: fmt.Sprintf("%s | %s | %s | %s",
			now.Format("15:04:05"), "POSITION: close", m.currentPosition.Side, comment),
	}); err != nil {
		m.log.Error(ctx, "error closing position marker",
			"error", err,
		)
	}
	m.currentPosition = nil

	return nil
}
