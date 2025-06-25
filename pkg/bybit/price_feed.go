package bybit

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/kirill-a-belov/trader/pkg/ticker"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

func (m *Module) PriceFeedBTCUSDT(ctx context.Context) (chan *ticker.Ticker, error) {
	_, span := tracer.Start(ctx, "pkg.bybit.PriceFeedBTCUSDT")
	defer span.End()

	result := make(chan *ticker.Ticker)

	go func() {
		const reconnectionDelay = time.Second
		for {
			err := m.wssTickerFeed(ctx, result)
			m.log.Error(ctx, "wssTickerFeed",
				"err", err,
				"next", "reconnecting",
				"delay", reconnectionDelay,
			)
			time.Sleep(reconnectionDelay)
		}
	}()

	return result, nil
}

func (m *Module) wssTickerFeed(ctx context.Context, result chan *ticker.Ticker) error {
	_, span := tracer.Start(ctx, "pkg.bybit.wssTickerFeed")
	defer span.End()

	const endpoint = "/v5/public/linear"
	conn, _, err := m.wssClient.Dial(m.config.ApiWssUrl+endpoint, nil)
	if err != nil {
		return errors.Wrap(err, "failed to dial wss")
	}
	defer conn.Close()
	m.wssPingPongHandler(ctx, conn)

	const (
		topicTickerName = "tickers.BTCUSDT"
		topicTradeName  = "publicTrade.BTCUSDT"
	)
	subscribe := map[string]interface{}{
		"op":   "subscribe",
		"args": []string{topicTickerName, topicTradeName},
	}
	if err := conn.WriteJSON(subscribe); err != nil {
		return errors.Wrap(err, "failed to send subscribe")
	}

	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			m.log.Error(ctx, "failed to read wss msg", "err", err)

			continue
		}

		var newTicker *ticker.Ticker
		switch msg["topic"] {
		case topicTickerName:
			newTicker, err = decodeTickerMessage(ctx, msg)
			if err != nil {
				m.log.Error(ctx, "failed to decode ticker message",
					"err", err,
					"msg", msg,
				)
				continue
			}
		case topicTradeName:
			newTicker, err = decodeTradeMessage(ctx, msg)
			if err != nil {
				m.log.Error(ctx, "failed to decode trade message",
					"err", err,
					"msg", msg,
				)
				continue
			}
		default:
		}

		if newTicker == nil {
			continue
		}

		fillEmptyFields(newTicker)

		result <- newTicker

	}
}

func (m *Module) wssPingPongHandler(ctx context.Context, conn *websocket.Conn) {
	_, span := tracer.Start(ctx, "pkg.bybit.wssPingPongHandler")
	defer span.End()

	conn.SetPongHandler(func(appData string) error {
		m.log.Debug(ctx, "wss PriceFeedBTCUSDT pong")

		return nil
	})
	go func() {
		const pingInterval = time.Second * 15

		ticker := time.NewTicker(pingInterval)
		defer ticker.Stop()
		for range ticker.C {
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				m.log.Error(ctx, "wss PriceFeedBTCUSDT ping error", err)
			}
			m.log.Debug(ctx, "wss PriceFeedBTCUSDT ping")
		}
	}()
}

func decodeTradeMessage(ctx context.Context, msg map[string]interface{}) (*ticker.Ticker, error) {
	_, span := tracer.Start(ctx, "pkg.bybit.decodeTradeMessage")
	defer span.End()

	result := &ticker.Ticker{}

	ts := time.UnixMilli(int64(msg["ts"].(float64)))
	if ts.IsZero() {
		return nil, errors.New("invalid timestamp")
	}
	result.Timestamp = ts

	data, ok := msg["data"].([]interface{})
	if !ok {
		return nil, errors.New("invalid data")
	}

	const (
		operationNameBuy  = "Buy"
		operationNameSell = "Sell"
	)
	for _, itemRaw := range data {
		item, ok := itemRaw.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid data item")
		}

		var (
			operation string
			volume    float64
			price     float64
			buyTick   bool
			typeTick  string
			timestamp time.Time
		)

		operationStr, ok := item["S"].(string)
		if ok && operationStr != "" {
			operation = operationStr
		}

		typeTickStr, ok := item["L"].(string)
		if ok && typeTickStr != "" {
			typeTick = typeTickStr
		}

		volumeStr, ok := item["v"].(string)
		if ok && volumeStr != "" {
			v, err := strconv.ParseFloat(volumeStr, 64)
			if err != nil {
				m.log.Error(ctx, "invalid volume in ticker")
			}
			volume = v
		}

		priceStr, ok := item["p"].(string)
		if ok && priceStr != "" {
			v, err := strconv.ParseFloat(priceStr, 64)
			if err != nil {
				m.log.Error(ctx, "invalid price in ticker")
			}
			price = v
		}

		buyTickStr, ok := item["BT"].(string)
		if ok && buyTickStr != "" {
			v, err := strconv.ParseBool(buyTickStr)
			if err != nil {
				m.log.Error(ctx, "invalid buyTick in ticker")
			}
			buyTick = v
		}

		timeStr := fmt.Sprint(item["T"])
		if timeStr != "" {
			v, err := strconv.ParseFloat(timeStr, 64)
			if err != nil {
				m.log.Error(ctx, "invalid timestamp in ticker")
			}
			timestamp = time.UnixMilli(int64(v))
		}

		switch operation {
		case operationNameBuy:
			result.BuyCount++
			result.BuyVolume += volume
		case operationNameSell:
			result.SellCount++
			result.SellVolume += volume
		}

		result.TradeList = append(result.TradeList, &ticker.Trade{
			Timestamp: timestamp,
			Price:     price,
			Volume:    volume,
			Side:      operation,
			BuyTick:   buyTick,
			TickType:  typeTick,
		})

	}

	return result, nil
}

func decodeTickerMessage(ctx context.Context, msg map[string]interface{}) (*ticker.Ticker, error) {
	_, span := tracer.Start(ctx, "pkg.bybit.decodeTickerMessage")
	defer span.End()

	result := &ticker.Ticker{}

	ts := time.UnixMilli(int64(msg["ts"].(float64)))
	if ts.IsZero() {
		return nil, errors.New("invalid timestamp")
	}
	result.Timestamp = ts

	data, ok := msg["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid data")
	}

	markPriceStr, ok := data["markPrice"].(string)
	if ok && markPriceStr != "" {
		markPrice, err := strconv.ParseFloat(markPriceStr, 64)
		if err != nil {
			m.log.Error(ctx, "invalid mark price in ticker")
		}
		result.MarkPrice = markPrice
	}

	indexPriceStr, ok := data["indexPrice"].(string)
	if ok && indexPriceStr != "" {
		indexPrice, err := strconv.ParseFloat(indexPriceStr, 64)
		if err != nil {
			m.log.Error(ctx, "invalid index price in ticker")
		}
		result.IndexPrice = indexPrice

	}

	lastPriceStr, ok := data["lastPrice"].(string)
	if ok && lastPriceStr != "" {
		lastPrice, err := strconv.ParseFloat(lastPriceStr, 64)
		if err != nil {
			m.log.Error(ctx, "invalid last price in ticker")
		}
		result.LastPrice = lastPrice
	}

	bid1PriceStr, ok := data["bid1Price"].(string)
	if ok && bid1PriceStr != "" {
		bid1Price, err := strconv.ParseFloat(bid1PriceStr, 64)
		if err != nil {
			m.log.Error(ctx, "invalid bid1 price in ticker")
		}
		result.Bid1Price = bid1Price
	}

	bid1SizeStr, ok := data["bid1Size"].(string)
	if ok && bid1SizeStr != "" {
		bid1Size, err := strconv.ParseFloat(bid1SizeStr, 64)
		if err != nil {
			m.log.Error(ctx, "invalid bid1 size in ticker")
		}
		result.Bid1Size = bid1Size
	}

	ask1PriceStr, ok := data["ask1Price"].(string)
	if ok && ask1PriceStr != "" {
		ask1Price, err := strconv.ParseFloat(ask1PriceStr, 64)
		if err != nil {
			m.log.Error(ctx, "invalid ask1 price in ticker")
		}
		result.Ask1Price = ask1Price
	}

	ask1SizeStr, ok := data["ask1Size"].(string)
	if ok && ask1SizeStr != "" {
		ask1Size, err := strconv.ParseFloat(ask1SizeStr, 64)
		if err != nil {
			m.log.Error(ctx, "invalid ask1 size in ticker")
		}
		result.Ask1Size = ask1Size
	}

	openInterestStr, ok := data["openInterest"].(string)
	if ok && openInterestStr != "" {
		openInterest, err := strconv.ParseFloat(openInterestStr, 64)
		if err != nil {
			m.log.Error(ctx, "invalid openInterest in ticker")
		}
		result.OpenInterest = openInterest
	}

	openInterestValueStr, ok := data["openInterestValue"].(string)
	if ok && openInterestValueStr != "" {
		openInterestValue, err := strconv.ParseFloat(openInterestValueStr, 64)
		if err != nil {
			m.log.Error(ctx, "invalid openInterestValue in ticker")
		}
		result.OpenInterestValue = openInterestValue
	}

	return result, nil
}

var lastTicker = &ticker.Ticker{}

func fillEmptyFields(newTicker *ticker.Ticker) {
	if newTicker.LastPrice == 0 {
		newTicker.LastPrice = lastTicker.LastPrice
	} else {
		lastTicker.LastPrice = newTicker.LastPrice
	}
	if newTicker.MarkPrice == 0 {
		newTicker.MarkPrice = lastTicker.MarkPrice
	} else {
		lastTicker.MarkPrice = newTicker.MarkPrice
	}
	if newTicker.IndexPrice == 0 {
		newTicker.IndexPrice = lastTicker.IndexPrice
	} else {
		lastTicker.IndexPrice = newTicker.IndexPrice
	}
	if newTicker.OpenInterest == 0 {
		newTicker.OpenInterest = lastTicker.OpenInterest
	} else {
		lastTicker.OpenInterest = newTicker.OpenInterest
	}
	if newTicker.OpenInterestValue == 0 {
		newTicker.OpenInterestValue = lastTicker.OpenInterestValue
	} else {
		lastTicker.OpenInterestValue = newTicker.OpenInterestValue
	}
	if newTicker.Bid1Size == 0 {
		newTicker.Bid1Size = lastTicker.Bid1Size
	} else {
		lastTicker.Bid1Size = newTicker.Bid1Size
	}
	if newTicker.Bid1Price == 0 {
		newTicker.Bid1Price = lastTicker.Bid1Price
	} else {
		lastTicker.Bid1Price = newTicker.Bid1Price
	}
	if newTicker.Ask1Size == 0 {
		newTicker.Ask1Size = lastTicker.Ask1Size
	} else {
		lastTicker.Ask1Size = newTicker.Ask1Size
	}
	if newTicker.Ask1Price == 0 {
		newTicker.Ask1Price = lastTicker.Ask1Price
	} else {
		lastTicker.Ask1Price = newTicker.Ask1Price
	}
}
