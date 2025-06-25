package chart

import (
	"context"
	"time"

	"github.com/kirill-a-belov/trader/internal/chart/model"
	signalModel "github.com/kirill-a-belov/trader/internal/signal/model"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

func (m *Module) PutCandle(ctx context.Context, candle *model.Candle) error {
	_, span := tracer.Start(ctx, "pkg.internal.chart.PutCandle")
	defer span.End()

	if len(m.candleList) > 0 {
		lastCandle := m.candleList[len(m.candleList)-1]
		if lastCandle != nil && lastCandle.Time == candle.Time {
			*lastCandle = *candle

			return nil
		}
	}

	m.candleList = append(m.candleList, candle)
	if len(m.candleList) > m.config.Depth {
		m.candleList = m.candleList[1:]
	}

	return nil
}

func (m *Module) PutMarker(ctx context.Context, marker *model.Marker) error {
	_, span := tracer.Start(ctx, "pkg.internal.chart.PutMarker")
	defer span.End()

	markerTime := marker.Timestamp.Truncate(m.config.Timeframe)
	marker.Time = markerTime.Unix()
	delete(m.markerStorage, markerTime.Add(-(time.Duration(m.config.Depth) * m.config.Timeframe)))

	_, ok := m.markerStorage[markerTime]
	if !ok {
		m.markerStorage[markerTime] = []*model.Marker{}
	}

	for i := range m.markerStorage[markerTime] {
		if (m.markerStorage[markerTime][i].Timeframe == marker.Timeframe) && (m.markerStorage[markerTime][i].Name == marker.Name) {
			m.markerStorage[markerTime][i] = marker

			return nil
		}
	}

	m.markerStorage[markerTime] = append(m.markerStorage[markerTime], marker)

	return nil
}

func (m *Module) FromSignal(s *signalModel.Signal) *model.Marker {
	var (
		colour string
		shape  string
	)
	switch s.Direction {
	case signalModel.SignalDirectionBull:
		colour = "green"
		shape = "▲"
	case signalModel.SignalDirectionBear:
		colour = "red"
		shape = "▼"
	case signalModel.SignalDirectionFlat:
		colour = "grey"
		shape = "◼"
	}

	marker := &model.Marker{}
	marker.Timestamp = s.Timestamp
	marker.Timeframe = s.Timeframe
	marker.Name = s.Name
	marker.Details = s.Description
	marker.Color = colour
	marker.Text = shape
	marker.Position = "aboveBar"

	return marker
}
