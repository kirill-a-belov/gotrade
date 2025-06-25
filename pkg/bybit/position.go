package bybit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/kirill-a-belov/trader/pkg/tracer"
)

type PositionBTCUSDTResponse struct {
	Symbol                 string `json:"symbol"`
	Leverage               string `json:"leverage"`
	AutoAddMargin          int    `json:"autoAddMargin"`
	AvgPrice               string `json:"avgPrice"`
	LiqPrice               string `json:"liqPrice"`
	RiskLimitValue         string `json:"riskLimitValue"`
	TakeProfit             string `json:"takeProfit"`
	PositionValue          string `json:"positionValue"`
	IsReduceOnly           bool   `json:"isReduceOnly"`
	TpslMode               string `json:"tpslMode"`
	RiskId                 int    `json:"riskId"`
	TrailingStop           string `json:"trailingStop"`
	UnrealisedPnl          string `json:"unrealisedPnl"`
	MarkPrice              string `json:"markPrice"`
	AdlRankIndicator       int    `json:"adlRankIndicator"`
	CumRealisedPnl         string `json:"cumRealisedPnl"`
	PositionMM             string `json:"positionMM"`
	CreatedTime            string `json:"createdTime"`
	PositionIdx            int    `json:"positionIdx"`
	PositionIM             string `json:"positionIM"`
	Seq                    int64  `json:"seq"`
	UpdatedTime            string `json:"updatedTime"`
	Side                   string `json:"side"`
	BustPrice              string `json:"bustPrice"`
	PositionBalance        string `json:"positionBalance"`
	LeverageSysUpdatedTime string `json:"leverageSysUpdatedTime"`
	CurRealisedPnl         string `json:"curRealisedPnl"`
	Size                   string `json:"size"`
	PositionStatus         string `json:"positionStatus"`
	MmrSysUpdatedTime      string `json:"mmrSysUpdatedTime"`
	StopLoss               string `json:"stopLoss"`
	TradeMode              int    `json:"tradeMode"`
	SessionAvgPrice        string `json:"sessionAvgPrice"`
}

func (m *Module) PositionBTCUSDT(ctx context.Context) (*PositionBTCUSDTResponse, error) {
	_, span := tracer.Start(ctx, "pkg.bybit.PositionBTCUSDT")
	defer span.End()

	request := m.restClient.R()
	query := "category=linear&symbol=BTCUSDT"

	m.setHeaders(ctx, request, query)

	const endpoint = "/v5/position/list"

	resp, err := request.Get(fmt.Sprintf("%s%s?%s", m.config.ApiRestUrl, endpoint, query))
	if err != nil {
		panic(err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.Errorf("invalid status: %s", resp.Status())
	}

	bbr := &bybitPositionListResponse{}
	if err := json.Unmarshal(resp.Body(), &bbr); err != nil {
		return nil, errors.Wrap(err, "unmarshalling response")
	}

	const btcusdtCoinName = "BTCUSDT"

	if len(bbr.Result.List) > 1 {
		return nil, errors.Errorf("multiple positions found: not supported")
	}
	if len(bbr.Result.List) < 1 {
		return nil, errors.Errorf("no positions found")
	}

	position := bbr.Result.List[0]

	if position.Symbol != btcusdtCoinName {
		return nil, errors.Errorf("invalid position symbol: %v", position)
	}

	return &position, nil
}

type bybitPositionListResponse struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		NextPageCursor string                    `json:"nextPageCursor"`
		Category       string                    `json:"category"`
		List           []PositionBTCUSDTResponse `json:"list"`
	} `json:"result"`
	RetExtInfo struct{} `json:"retExtInfo"`
	Time       int64    `json:"time"`
}
