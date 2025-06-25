package simple_algo

import (
	"context"
	"fmt"

	positionModel "github.com/kirill-a-belov/trader/internal/position/model"
	"github.com/kirill-a-belov/trader/pkg/ticker"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

func (m *Module) shouldClosePositionByRiskManagement(
	ctx context.Context,
	priceTicker *ticker.Ticker,
	position *positionModel.Position,
) (
	bool,
	string, error,
) {
	_, span := tracer.Start(ctx, "pkg.internal.simple_algo.shouldClosePositionByRiskManagement")
	defer span.End()

	clearPnl := m.clearPNL(ctx, priceTicker, position)

	if position.Age() > m.config.NonProfitPositionMaxAge {
		if clearPnl > 0.05 {
			if m.config.Debug {
				fmt.Println("closing by age and small profit",
					"clearPnl", clearPnl,
				)
			}

			return true, "closing by age and small profit", nil
		}
	}

	nextCandleBegins := priceTicker.Timestamp.Second() > 4 && priceTicker.Timestamp.Second() > 8
	if nextCandleBegins {
		const riskSizeUSDT = 1
		if clearPnl <= -riskSizeUSDT {
			if m.config.Debug {
				fmt.Println("closing by next candle begins and max risk loss",
					"clearPnl", clearPnl,
				)
			}

			return true, "closing by next candle begins and max risk loss", nil
		}
	}

	return false, "", nil
}
