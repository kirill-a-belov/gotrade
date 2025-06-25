package model

import (
	"fmt"
	"math"
	"time"

	"github.com/kirill-a-belov/trader/pkg/ticker"
)

type Candle struct {
	Timestamp time.Time

	High  float64
	Low   float64
	Open  float64
	Close float64

	TickCount float64
	TickSum   float64

	SellCount  float64
	SellVolume float64
	BuyCount   float64
	BuyVolume  float64

	TradeList []*ticker.Trade
}

func (c *Candle) Print() {
	const (
		colorReset = "\033[0m"
		colorRed   = "\033[31m"
		colorGreen = "\033[32m"
	)

	var direction string
	var color string

	switch {
	case c.Bull():
		direction = "▲ Bull"
		color = colorGreen
	case c.Bear():
		direction = "▼ Bear"
		color = colorRed
	default:
		direction = "◆ Neutral"
		color = colorReset
	}
	fmt.Printf("Open -> Close:			%0.2f -> %0.2f\n", c.Open, c.Close)
	fmt.Printf("Direction | Diff | Diff%%:	%s%-6s%s | %0.2f | %0.2f\n", color, direction, colorReset, c.Close-c.Open, (c.Close-c.Open)*100/c.Close)
	fmt.Printf("Low-> High:			%0.2f -> %0.2f\n", c.Low, c.High)
	fmt.Printf("Ticks: Count | Sum | AvgPrice:	%0.0f | %0.2f | %0.2f \n", c.TickCount, c.TickSum, c.Average())
	fmt.Printf("Trades: (count|volume|avg):	Buy(%0.0f|%0.3f|%0.3f), Sell(%0.0f|%0.3f|%0.3f), Total(%0.0f|%0.3f|%0.3f)\n",
		c.BuyCount, c.BuyVolume, c.BuyVolume/c.BuyCount,
		c.SellCount, c.SellVolume, c.SellVolume/c.SellCount,
		c.BuyCount+c.SellCount, c.BuyVolume+c.SellVolume, (c.BuyVolume+c.SellVolume)/(c.BuyCount+c.SellCount),
	)
}

func (c *Candle) Bull() bool {
	return c.Close > c.Open
}

func (c *Candle) Bear() bool {
	return c.Close < c.Open
}

func (c *Candle) Average() float64 {
	if c.TickCount == 0 {
		return 0
	}
	return c.TickSum / c.TickCount
}

func (c *Candle) BodySize() float64 {
	return math.Abs(c.Close - c.Open)
}

func (c *Candle) ShadowSize() float64 {
	return (c.High - c.Low) - c.BodySize()
}

func (c *Candle) Volatility() float64 {
	return c.High - c.Low
}

func (c *Candle) Momentum() float64 {
	return c.Close - c.Open
}

func (c *Candle) IsDoji(threshold float64) bool {
	return math.Abs(c.Close-c.Open) < threshold*(c.High-c.Low)
}

func (c *Candle) IsImpulse(threshold float64) bool {
	return c.BodySize() > threshold*(c.High-c.Low)
}
