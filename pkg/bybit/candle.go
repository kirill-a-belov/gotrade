package bybit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/kirill-a-belov/trader/pkg/tracer"
)

type CandleListBTCUSDTResponse []*Kline

func (m *Module) CandleListBTCUSDT(ctx context.Context) (CandleListBTCUSDTResponse, error) {
	_, span := tracer.Start(ctx, "pkg.bybit.CandleListBTCUSDT")
	defer span.End()

	request := m.restClient.R()
	query := "category=linear&symbol=BTCUSDT&interval=1&limit=100"

	m.setHeaders(ctx, request, query)

	const endpoint = "/v5/market/kline"

	resp, err := request.Get(fmt.Sprintf("%s%s?%s", m.config.ApiRestUrl, endpoint, query))
	if err != nil {
		panic(err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.Errorf("invalid status: %s", resp.Status())
	}

	bkr := &bybitKlineResponse{}
	if err := json.Unmarshal(resp.Body(), bkr); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response body")
	}

	result := make(CandleListBTCUSDTResponse, len(bkr.Result.List))
	for i := range bkr.Result.List {
		result[i], err = parseKline(bkr.Result.List[i])
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse kline")
		}
	}

	return result, nil
}

const KlineTimeframe = time.Minute

type Kline struct {
	StartTime time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	Turnover  float64
}

func parseKline(record []string) (*Kline, error) {
	if len(record) < 7 {
		return nil, fmt.Errorf("invalid record length: %d", len(record))
	}

	startTime, err := strconv.ParseInt(record[0], 10, 64)
	if err != nil {
		return nil, err
	}

	open, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return nil, err
	}

	high, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return nil, err
	}

	low, err := strconv.ParseFloat(record[3], 64)
	if err != nil {
		return nil, err
	}

	closePrice, err := strconv.ParseFloat(record[4], 64)
	if err != nil {
		return nil, err
	}

	volume, err := strconv.ParseFloat(record[5], 64)
	if err != nil {
		return nil, err
	}

	turnover, err := strconv.ParseFloat(record[6], 64)
	if err != nil {
		return nil, err
	}

	return &Kline{
		StartTime: time.UnixMilli(startTime),
		Open:      open,
		High:      high,
		Low:       low,
		Close:     closePrice,
		Volume:    volume,
		Turnover:  turnover,
	}, nil
}

type bybitKlineResponse struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		Symbol   string     `json:"symbol"`
		Category string     `json:"category"`
		List     [][]string `json:"list"`
	} `json:"result"`
	RetExtInfo struct{} `json:"retExtInfo"`
	Time       int64    `json:"time"`
}
