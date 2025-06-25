package model

import (
	"fmt"
	"time"
)

const (
	SignalDirectionBull = "bull"
	SignalDirectionBear = "bear"
	SignalDirectionFlat = "flat"
)

type Signal struct {
	Timestamp   time.Time
	Timeframe   time.Duration
	Description string
	Direction   string
	Name        string
	Confidence  float64
}

func TimeframeAndName(timeframe time.Duration, name string) string {
	return fmt.Sprint(timeframe, "+", name)
}

type SignalList []*Signal

func (sl *SignalList) TimeframeAndNameMap() map[string]*Signal {
	result := make(map[string]*Signal)

	for _, signal := range *sl {
		result[TimeframeAndName(signal.Timeframe, signal.Name)] = signal
	}

	return result
}
