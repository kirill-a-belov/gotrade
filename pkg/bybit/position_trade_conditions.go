package bybit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/kirill-a-belov/trader/pkg/tracer"
)

type SetTradeConditionsBTCUSDTRequest struct {
	Category    string `json:"category"`    // Required: Product type, e.g. "linear", "inverse"
	Symbol      string `json:"symbol"`      // Required: Symbol name in uppercase, e.g. "BTCUSDT"
	TpslMode    string `json:"tpslMode"`    // Required: TP/SL mode ("Full" or "Partial")
	PositionIdx int    `json:"positionIdx"` // Required: Position index (0 = one-way, 1 = hedge buy, 2 = hedge sell)

	TakeProfit   *string `json:"takeProfit,omitempty"`   // Optional: Take profit price, >= 0, 0 to cancel TP
	StopLoss     *string `json:"stopLoss,omitempty"`     // Optional: Stop loss price, >= 0, 0 to cancel SL
	TrailingStop *string `json:"trailingStop,omitempty"` // Optional: Trailing stop distance, >= 0, 0 to cancel trailing stop

	TpTriggerBy *string `json:"tpTriggerBy,omitempty"` // Optional: TP trigger price type ("MarkPrice", "IndexPrice", "LastPrice")
	SlTriggerBy *string `json:"slTriggerBy,omitempty"` // Optional: SL trigger price type ("MarkPrice", "IndexPrice", "LastPrice")

	ActivePrice *string `json:"activePrice,omitempty"` // Optional: Activation price for trailing stop (trigger TS only after this price)

	TpSize *string `json:"tpSize,omitempty"` // Optional: TP size for partial mode; must equal slSize if used
	SlSize *string `json:"slSize,omitempty"` // Optional: SL size for partial mode; must equal tpSize if used

	TpLimitPrice *string `json:"tpLimitPrice,omitempty"` // Optional: Limit price for TP order (only for Partial mode + tpOrderType=Limit)
	SlLimitPrice *string `json:"slLimitPrice,omitempty"` // Optional: Limit price for SL order (only for Partial mode + slOrderType=Limit)

	TpOrderType *string `json:"tpOrderType,omitempty"` // Optional: TP order type ("Market" or "Limit"); Full mode supports only "Market"
	SlOrderType *string `json:"slOrderType,omitempty"` // Optional: SL order type ("Market" or "Limit"); Full mode supports only "Market"
}

func (m *Module) SetTradeConditionsBTCUSDT(ctx context.Context, request *SetTradeConditionsBTCUSDTRequest) error {
	_, span := tracer.Start(ctx, "pkg.bybit.SetTradeConditionsBTCUSDT")
	defer span.End()

	client := m.restClient.R()

	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return errors.Wrap(err, "failed to marshal request body")
	}

	m.setHeaders(ctx, client, string(bodyBytes))
	client.SetBody(bodyBytes)

	const endpoint = "/v5/position/trading-stop"

	resp, err := client.Post(fmt.Sprintf("%s%s", m.config.ApiRestUrl, endpoint))
	if err != nil {
		panic(err)
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.Errorf("invalid status: %s", resp.Status())
	}

	bbr := &bybitTradingStopResponse{}
	if err := json.Unmarshal(resp.Body(), &bbr); err != nil {
		return errors.Wrap(err, "unmarshalling response")
	}

	const ok = "OK"
	if bbr.RetMsg != ok {
		return errors.Errorf("invalid response: %v", bbr)
	}

	return nil
}

type bybitTradingStopResponse struct {
	RetCode    int      `json:"retCode"`
	RetMsg     string   `json:"retMsg"`
	Result     struct{} `json:"result"`
	RetExtInfo struct{} `json:"retExtInfo"`
	Time       int64    `json:"time"`
}
