package candle

import (
	"context"
	"github.com/kirill-a-belov/trader/pkg/tracer"
	"time"

	"github.com/kirill-a-belov/trader/internal/candle/model"
)

func (m *Module) candleListOldToNew(ctx context.Context, timestamp time.Time, period int) []*model.Candle {
	span, _ := tracer.Start(ctx, "pkg.internal.candle.candleListOldToNew")
	defer span.Done()

	result := make([]*model.Candle, 0, period)
	for i := period - 1; i >= 0; i-- {
		candle := m.Candle(ctx, timestamp.Add(-m.Config.Timeframe*time.Duration(i)))
		if candle == nil {
			return nil
		}
		result = append(result, candle)
	}

	return result
}
