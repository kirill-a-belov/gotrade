package model

import "time"

type Position struct {
	Id            string
	CreatedAt     time.Time
	Side          PositionSide
	Size          float64
	AvgPrice      float64
	StopLossPrice float64
}

type PositionSide string

const (
	PositionSideUnknown PositionSide = "Unknown"
	PositionSideBuy     PositionSide = "Buy"
	PositionSideSell    PositionSide = "Sell"
)

func (p *Position) Age() time.Duration {
	return time.Since(p.CreatedAt)
}

func (p *Position) OppositeSide() PositionSide {
	switch p.Side {
	case PositionSideBuy:
		return PositionSideSell
	case PositionSideSell:
		return PositionSideBuy
	default:
		return PositionSideUnknown
	}
}
