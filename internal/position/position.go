package position

import (
	"context"
	"strconv"

	"github.com/pkg/errors"

	"github.com/kirill-a-belov/trader/internal/position/model"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

func (m *Module) Position(ctx context.Context) (*model.Position, error) {
	_, span := tracer.Start(ctx, "pkg.internal.position.Position")
	defer span.End()

	if m.currentPosition == nil {
		return nil, nil
	}

	// TODO Fetch position with wss
	position, err := m.bybitModule.PositionBTCUSDT(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "bybitModule.PositionBTCUSDT")
	}
	if position.Size == "0" {
		m.log.Error(ctx, "position doesn't exists")
		m.currentPosition = nil

		return nil, nil
	}

	positionSize, err := strconv.ParseFloat(position.Size, 64)
	if err != nil {
		return nil, errors.Wrap(err, "strconv.ParseFloat position size")
	}
	if m.currentPosition.Size != positionSize {
		m.log.Error(ctx, "position size does not match")
		m.currentPosition.Size = positionSize
	}

	positionStopLossPrice, err := strconv.ParseFloat(position.StopLoss, 64)
	if err != nil {
		return nil, errors.Wrap(err, "strconv.ParseFloat position stop loss price")
	}
	if m.currentPosition.StopLossPrice != positionStopLossPrice {
		m.log.Error(ctx, "position stop loss price does not match")
		m.currentPosition.StopLossPrice = positionStopLossPrice
	}

	if string(m.currentPosition.Side) != position.Side {
		m.log.Error(ctx, "position side does not match")
		m.currentPosition.Side = model.PositionSide(position.Side)
	}

	positionAvgPrice, err := strconv.ParseFloat(position.AvgPrice, 64)
	if err != nil {
		return nil, errors.Wrap(err, "strconv.ParseFloat position avg price")
	}
	m.currentPosition.AvgPrice = positionAvgPrice

	return m.currentPosition, nil
}
