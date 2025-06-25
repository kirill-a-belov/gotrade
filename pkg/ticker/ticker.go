package ticker

import "time"

type Ticker struct {
	Timestamp  time.Time `json:"timestamp"`
	MarkPrice  float64   `json:"markPrice"`
	IndexPrice float64   `json:"indexPrice"`
	LastPrice  float64   `json:"lastPrice"`

	Bid1Price float64 `json:"bid1Price"`
	Bid1Size  float64 `json:"bid1Size"`

	Ask1Price float64 `json:"ask1Price"`
	Ask1Size  float64 `json:"ask1Size"`

	OpenInterest      float64 `json:"openInterest"`
	OpenInterestValue float64 `json:"openInterestValue"`

	SellCount  float64 `json:"SellCount"`
	SellVolume float64 `json:"sellVolume"`
	BuyCount   float64 `json:"BuyCount"`
	BuyVolume  float64 `json:"buyVolume"`

	TradeList []*Trade `json:"tradeList"`
}

func (t *Ticker) IsTooOld() bool {
	return time.Since(t.Timestamp) > 3*time.Second
}

type Trade struct {
	Timestamp time.Time `json:"timestamp"`
	Price     float64   `json:"price"`
	Volume    float64   `json:"volume"`
	Side      string    `json:"side"`
	BuyTick   bool      `json:"buyTick"`
	TickType  string    `json:"tickType"`
}
