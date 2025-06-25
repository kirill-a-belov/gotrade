package simple_algo

import (
	"context"

	"github.com/pkg/errors"

	positionModel "github.com/kirill-a-belov/trader/internal/position/model"
	signalModel "github.com/kirill-a-belov/trader/internal/signal/model"
	"github.com/kirill-a-belov/trader/pkg/ticker"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

func (m *Module) clearPNL(ctx context.Context, priceTicker *ticker.Ticker, position *positionModel.Position) float64 {
	_, span := tracer.Start(ctx, "pkg.internal.simple_algo.clearPNL")
	defer span.End()

	var pnl float64
	switch position.Side {
	case positionModel.PositionSideSell:
		pnl = (position.AvgPrice - priceTicker.LastPrice) * position.Size
	case positionModel.PositionSideBuy:
		pnl = (priceTicker.LastPrice - position.AvgPrice) * position.Size
	}

	unrealizedPnlWithCommissionAndSlippage := pnl - m.config.CommissionInAndOutSizeUSDT - (pnl * m.config.SlippagePercent / 100)
	clearPnl := unrealizedPnlWithCommissionAndSlippage

	return clearPnl
}

func (m *Module) managePosition(ctx context.Context, priceTicker *ticker.Ticker, sl signalModel.SignalList) error {
	_, span := tracer.Start(ctx, "pkg.internal.simple_algo.managePosition")
	defer span.End()

	currentPosition, err := m.positionModule.Position(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting current position")
	}
	if currentPosition == nil {
		if priceTicker.Timestamp.Second() > m.config.MaxSecondsForEnter {
			return nil
		}

		newPositionSide, comment, err := m.canOpenPositionBySignalList(ctx, sl)
		if err != nil {
			return errors.Wrap(err, "error opening position side")
		}
		if newPositionSide == positionModel.PositionSideUnknown {
			return nil
		}
		if err := m.positionModule.OpenPosition(ctx, priceTicker.LastPrice, newPositionSide, comment); err != nil {
			return errors.Wrap(err, "error opening position side")
		}

		return nil
	}

	shouldCloseByRiskManagement, comment, err := m.shouldClosePositionByRiskManagement(ctx, priceTicker, currentPosition)
	if err != nil {
		return errors.Wrap(err, "error checking if position is closed by risk management")
	}
	if shouldCloseByRiskManagement {
		err = m.positionModule.ClosePosition(ctx, comment)
		if err != nil {
			return errors.Wrap(err, "error closing position")
		}

		return nil
	}

	shouldCloseBySignalList, comment, err := m.shouldClosePositionBySignalLists(ctx, sl, priceTicker, currentPosition)
	if err != nil {
		return errors.Wrap(err, "error checking if position is closed by signal list")
	}
	if shouldCloseBySignalList {
		err = m.positionModule.ClosePosition(ctx, comment)
		if err != nil {
			return errors.Wrap(err, "error closing position")
		}

		return nil
	}

	if err := m.positionModule.UpdatePositionStopLoss(ctx, priceTicker.LastPrice); err != nil {
		return errors.Wrap(err, "error updating position stop loss price")
	}

	// TODO Adjust position - more cases

	return nil
}
