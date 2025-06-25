package model

import (
	"github.com/kirill-a-belov/trader/pkg/ticker"
)

type Candle struct {
	Time  int64   `json:"time"`
	Open  float64 `json:"open"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Close float64 `json:"close"`

	BuyVolume  float64 `json:"buyVolume"`
	SellVolume float64 `json:"sellVolume"`
	BuyCount   float64 `json:"buyCount"`
	SellCount  float64 `json:"sellCount"`

	TickSum   float64 `json:"tickSum"`
	TickCount float64 `json:"tickCount"`

	TradeList []*ticker.Trade `json:"tradeList"`
}
