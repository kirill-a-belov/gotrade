package bybit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/kirill-a-belov/trader/pkg/tracer"
)

type CreateOrderRequest struct {
	Category       string  `json:"category"`                 // Market category: "linear", "inverse", etc.
	Symbol         string  `json:"symbol"`                   // Trading pair symbol, e.g. "BTCUSDT"
	Side           string  `json:"side"`                     // Order side: "Buy" or "Sell"
	OrderType      string  `json:"orderType"`                // Order type: "Limit", "Market", "Stop", "StopLimit", etc.
	Qty            string  `json:"qty"`                      // Order quantity (lots), string format
	Price          *string `json:"price,omitempty"`          // Order price, required for Limit and StopLimit orders
	TimeInForce    string  `json:"timeInForce"`              // Time in force: "GoodTillCancel", "ImmediateOrCancel", "FillOrKill"
	ReduceOnly     *bool   `json:"reduceOnly,omitempty"`     // Reduce-only order flag (optional)
	CloseOnTrigger *bool   `json:"closeOnTrigger,omitempty"` // Close position on trigger flag (optional)
	PositionIdx    *int    `json:"positionIdx,omitempty"`    // Position index for hedging mode (0,1,2) (optional)
	TriggerPrice   *string `json:"triggerPrice,omitempty"`   // Trigger price for stop orders (optional)
	TriggerBy      *string `json:"triggerBy,omitempty"`      // Price type to trigger by: "LastPrice", "IndexPrice", "MarkPrice" (optional)
	TakeProfit     *string `json:"takeProfit,omitempty"`     // Take profit price (optional)
	StopLoss       *string `json:"stopLoss,omitempty"`       // Stop loss price (optional)
	TpTriggerBy    *string `json:"tpTriggerBy,omitempty"`    // Price type for take profit trigger (optional)
	SlTriggerBy    *string `json:"slTriggerBy,omitempty"`    // Price type for stop loss trigger (optional)
	OrderLinkId    *string `json:"orderLinkId,omitempty"`    // Custom order ID for client (optional)
	BasePrice      *string `json:"basePrice,omitempty"`      // Base price for stop orders (optional)
}

func (m *Module) CreateOrder(ctx context.Context, request *CreateOrderRequest) (string, error) {
	_, span := tracer.Start(ctx, "pkg.bybit.CreateOrder")
	defer span.End()

	client := m.restClient.R()

	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal request body")
	}

	m.setHeaders(ctx, client, string(bodyBytes))
	client.SetBody(bodyBytes)

	const endpoint = "/v5/order/create"

	resp, err := client.Post(fmt.Sprintf("%s%s", m.config.ApiRestUrl, endpoint))
	if err != nil {
		panic(err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", errors.Errorf("invalid status: %s", resp.Status())
	}

	bbr := &bybitCreateOrderResponse{}
	if err := json.Unmarshal(resp.Body(), &bbr); err != nil {
		return "", errors.Wrap(err, "unmarshalling response")
	}

	const ok = "OK"
	if bbr.RetMsg != ok {
		return "", errors.Errorf("invalid response: %v", bbr)
	}

	return bbr.Result.OrderId, nil
}

type bybitCreateOrderResponse struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		OrderId     string `json:"orderId"`
		OrderLinkId string `json:"orderLinkId"`
	} `json:"result"`
	RetExtInfo struct{} `json:"retExtInfo"`
	Time       int64    `json:"time"`
}
