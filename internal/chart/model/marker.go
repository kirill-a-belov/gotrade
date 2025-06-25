package model

import "time"

type Marker struct {
	Timestamp time.Time
	Timeframe time.Duration
	Name      string

	Time     int64  `json:"time"`
	Text     string `json:"text"`
	Color    string `json:"color"`
	Position string `json:"position"` //"aboveBar" | "belowBar" | "inBar" | "atPriceTop" | "atPriceBottom" | "atPriceMiddle"
	Shape    string `json:"shape"`    //"circle" | "square" | "arrowUp" | "arrowDown"
	Details  string `json:"details"`
}
