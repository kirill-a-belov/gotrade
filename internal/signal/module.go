package signal

import (
	"context"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"github.com/kirill-a-belov/trader/pkg/logger"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

type Config struct {
	FastEmaPeriod int `envconfig:"TRADER_SIGNAL_FAST_EMA_PERIOD" default:"3"`
	SlowEmaPeriod int `envconfig:"TRADER_SIGNAL_SLOW_EMA_PERIOD" default:"6"`
	PricePeriod   int `envconfig:"TRADER_PRICE_PERIOD" default:"3"`

	PriceHugeCandlePercent int `envconfig:"TRADER_PRICE_HUGE_CANDLE_PERCENT" default:"1"`

	ExtraPeriod int           `envconfig:"TRADER_SIGNAL_EXTRA_PERIOD" default:"2"`
	Timeframe   time.Duration `envconfig:"TRADER_SIGNAL_TIMEFRAME" default:"1m"`
}

func (c *Config) Load() error {
	return envconfig.Process("", c)
}

func New(ctx context.Context) (*Module, error) {
	_, span := tracer.Start(ctx, "pkg.internal.signal.New")
	defer span.End()

	l := logger.New("signal")

	c := &Config{}
	if err := c.Load(); err != nil {
		return nil, errors.Wrap(err, "loading configuration")
	}

	m := &Module{
		Config: c,
		log:    l,
	}

	return m, nil
}

type Module struct {
	Config *Config
	log    logger.Logger
}

func (m *Module) CandleListNecessaryDepth() int {
	maxPeriod := 0
	for _, item := range []int{m.Config.FastEmaPeriod, m.Config.SlowEmaPeriod, m.Config.PricePeriod} {
		if item > maxPeriod {
			maxPeriod = item
		}
	}

	return maxPeriod + m.Config.ExtraPeriod
}
