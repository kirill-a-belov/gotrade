package bybit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	"github.com/kirill-a-belov/trader/pkg/tracer"
)

func (m *Module) EquityUSDT(ctx context.Context) (float64, error) {
	_, span := tracer.Start(ctx, "pkg.bybit.EquityUSDT")
	defer span.End()

	request := m.restClient.R()
	query := "accountType=UNIFIED"

	m.setHeaders(ctx, request, query)

	const endpoint = "/v5/account/wallet-balance"

	resp, err := request.Get(fmt.Sprintf("%s%s?%s", m.config.ApiRestUrl, endpoint, query))
	if err != nil {
		panic(err)
	}

	if resp.StatusCode() != http.StatusOK {
		return 0, errors.Errorf("invalid status: %s", resp.Status())
	}

	bbr := &bybitBalanceResponse{}
	if err := json.Unmarshal(resp.Body(), &bbr); err != nil {
		return 0, errors.Wrap(err, "unmarshalling response")
	}

	const usdtCoinName = "USDT"

	if len(bbr.Result.AccountList) > 1 {
		return 0, errors.Errorf("multiple USDT accounts found: not supported")
	}
	account := bbr.Result.AccountList[0]

	for _, coin := range account.CoinList {
		if coin.CoinName == usdtCoinName {
			equity, err := strconv.ParseFloat(coin.Equity, 64)
			if err != nil {
				return 0, errors.Wrapf(err, "invalid equity: %s", coin.Equity)
			}

			return equity, nil
		}
	}

	return 0, nil
}

type bybitBalanceResponse struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		AccountList []struct {
			TotalEquity            string `json:"totalEquity"`
			AccountIMRate          string `json:"accountIMRate"`
			TotalMarginBalance     string `json:"totalMarginBalance"`
			TotalInitialMargin     string `json:"totalInitialMargin"`
			AccountType            string `json:"accountType"`
			TotalAvailableBalance  string `json:"totalAvailableBalance"`
			AccountMMRate          string `json:"accountMMRate"`
			TotalPerpUPL           string `json:"totalPerpUPL"`
			TotalWalletBalance     string `json:"totalWalletBalance"`
			AccountLTV             string `json:"accountLTV"`
			TotalMaintenanceMargin string `json:"totalMaintenanceMargin"`
			CoinList               []struct {
				AvailableToBorrow   string `json:"availableToBorrow"`
				Bonus               string `json:"bonus"`
				AccruedInterest     string `json:"accruedInterest"`
				AvailableToWithdraw string `json:"availableToWithdraw"`
				TotalOrderIM        string `json:"totalOrderIM"`
				Equity              string `json:"equity"`
				TotalPositionMM     string `json:"totalPositionMM"`
				UsdValue            string `json:"usdValue"`
				UnrealisedPnl       string `json:"unrealisedPnl"`
				CollateralSwitch    bool   `json:"collateralSwitch"`
				SpotHedgingQty      string `json:"spotHedgingQty"`
				BorrowAmount        string `json:"borrowAmount"`
				TotalPositionIM     string `json:"totalPositionIM"`
				WalletBalance       string `json:"walletBalance"`
				CumRealisedPnl      string `json:"cumRealisedPnl"`
				Locked              string `json:"locked"`
				MarginCollateral    bool   `json:"marginCollateral"`
				CoinName            string `json:"coin"`
			} `json:"coin"`
		} `json:"list"`
	} `json:"result"`
	RetExtInfo struct{} `json:"retExtInfo"`
	Time       int64    `json:"time"`
}
