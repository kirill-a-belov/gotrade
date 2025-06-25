package candle

import (
	"context"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"github.com/kirill-a-belov/trader/internal/candle/model"
	"github.com/kirill-a-belov/trader/internal/chart"
	"github.com/kirill-a-belov/trader/internal/signal"
	"github.com/kirill-a-belov/trader/pkg/logger"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

type Config struct {
	Depth     time.Duration `envconfig:"TRADER_CANDLE_DEPTH" default:"60m"`
	Timeframe time.Duration `envconfig:"TRADER_CANDLE_TIMEFRAME" default:"1m"`

	SendCandleToChart bool `envconfig:"TRADER_CANDLE_SEND_CANDLE_TO_CHART" default:"true"`
	SendMarkerToChart bool `envconfig:"TRADER_CANDLE_SEND_MARKER_TO_CHART" default:"true"`
}

func (c *Config) Load() error {
	return envconfig.Process("", c)
}

func New(ctx context.Context, c *Config) (*Module, error) {
	_, span := tracer.Start(ctx, "pkg.internal.candle.New")
	defer span.End()

	l := logger.New("candle")

	if c == nil {
		c = &Config{}
		if err := c.Load(); err != nil {
			return nil, errors.Wrap(err, "loading configuration")
		}
	}

	cm, err := chart.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "creating chart model")
	}

	sm, err := signal.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "creating signal model")
	}
	sm.Config.Timeframe = c.Timeframe

	m := &Module{
		Config: c,
		log:    l,

		candleStorage: make(map[time.Time]*model.Candle),

		chartModule:  cm,
		signalModule: sm,
	}

	return m, nil
}

type Module struct {
	Config *Config
	log    logger.Logger

	candleStorage map[time.Time]*model.Candle
	chartModule   *chart.Module
	signalModule  *signal.Module
}
