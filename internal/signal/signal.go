package signal

import (
	"context"
	"fmt"
	"github.com/kirill-a-belov/trader/pkg/tracer"
	"time"

	"github.com/kirill-a-belov/trader/internal/signal/model"
)

func (m *Module) signal(
	ctx context.Context,
	name string,
	direction string,
	confidence float64,
) *model.Signal {
	_, span := tracer.Start(ctx, "pkg.internal.signal.signal")
	defer span.End()

	return &model.Signal{
		Timestamp: time.Now(),
		Timeframe: m.Config.Timeframe,
		Name:      name,
		Description: fmt.Sprintf("%s | %s | %s | %0.2f",
			time.Now().Format("15:04:05"), m.Config.Timeframe.String(), name, confidence),
		Direction:  direction,
		Confidence: confidence,
	}
}
