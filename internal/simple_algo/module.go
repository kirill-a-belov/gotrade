package simple_algo

import (
	"context"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"github.com/kirill-a-belov/trader/internal/candle"
	"github.com/kirill-a-belov/trader/internal/position"
	"github.com/kirill-a-belov/trader/pkg/bybit"
	"github.com/kirill-a-belov/trader/pkg/logger"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

type config struct {
	RiskSizeUSDT               float64       `envconfig:"SIMPLE_ALGO_RISK_SIZE_USDT" default:"1"`
	CommissionInAndOutSizeUSDT float64       `envconfig:"SIMPLE_ALGO_COMMISSION_IN_AND_OUT_SIZE_USDT" default:"0.09"`
	SlippagePercent            float64       `envconfig:"SIMPLE_ALGO_SLIPPAGE_PERCENT" default:"20"`
	NonProfitPositionMaxAge    time.Duration `envconfig:"SIMPLE_ALGO_NON_PROFIT_POSITION_MAX_AGE" default:"10m"`

	MaxSecondsForEnter int `envconfig:"SIMPLE_ALGO_MAX_SECONDS_FOR_ENTER" default:"5"`

	Debug bool `envconfig:"SIMPLE_ALGO_DEBUG" default:"false"`
}

func (c *config) Load() error {
	return envconfig.Process("", c)
}

var m *Module

func New(ctx context.Context) (*Module, error) {
	_, span := tracer.Start(ctx, "pkg.internal.simple_algo.New")
	defer span.End()

	if m != nil {
		return m, nil
	}

	l := logger.New("simple_algo")

	c := &config{}
	if err := c.Load(); err != nil {
		return nil, errors.Wrap(err, "loading configuration")
	}

	pm, err := position.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "loading position module")
	}
	bm, err := bybit.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "creating bybit module")
	}

	cnm1m, err := candle.New(ctx, &candle.Config{
		Depth:             60 * time.Minute,
		Timeframe:         time.Minute,
		SendCandleToChart: true,
		SendMarkerToChart: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "creating candle module")
	}

	cnm10, err := candle.New(ctx, &candle.Config{
		Depth:             10 * time.Minute,
		Timeframe:         10 * time.Second,
		SendCandleToChart: false,
		SendMarkerToChart: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "creating candle module 10")
	}

	cnm15, err := candle.New(ctx, &candle.Config{
		Depth:             15 * time.Minute,
		Timeframe:         15 * time.Second,
		SendCandleToChart: false,
		SendMarkerToChart: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "creating candle module 15")
	}
	cnm10m, err := candle.New(ctx, &candle.Config{
		Depth:             120 * time.Minute,
		Timeframe:         10 * time.Minute,
		SendCandleToChart: false,
		SendMarkerToChart: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "creating candle module 10m")
	}
	cnm5m, err := candle.New(ctx, &candle.Config{
		Depth:             240 * time.Minute,
		Timeframe:         5 * time.Minute,
		SendCandleToChart: false,
		SendMarkerToChart: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "creating candle module 5m")
	}
	candleModuleList := []*candle.Module{cnm10, cnm15, cnm1m, cnm5m, cnm10m}

	m = &Module{
		config: c,
		log:    l,

		candleModuleList: candleModuleList,

		positionModule: pm,
		bybitModule:    bm,
	}

	return m, nil
}

type Module struct {
	config *config
	log    logger.Logger

	candleModuleList []*candle.Module

	positionModule *position.Module

	bybitModule *bybit.Module
}

func (m *Module) Name() string {
	return "simple_trading_algo"
}
