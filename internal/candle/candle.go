package candle

import (
	"context"
	"time"

	"github.com/kirill-a-belov/trader/internal/candle/model"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

func (m *Module) Candle(ctx context.Context, ts time.Time) *model.Candle {
	_, span := tracer.Start(ctx, "pkg.internal.candle.Candle")
	defer span.End()

	candleTime := ts.Truncate(m.Config.Timeframe)

	candle, ok := m.candleStorage[candleTime]
	if !ok {
		return nil
	}

	return candle
}
