package position

import (
	"context"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"github.com/kirill-a-belov/trader/internal/chart"
	"github.com/kirill-a-belov/trader/internal/position/model"
	"github.com/kirill-a-belov/trader/pkg/bybit"
	"github.com/kirill-a-belov/trader/pkg/logger"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

type config struct {
	MaxRiskSizeUSDT float64 `envconfig:"POSITION_MAX_RISK_SIZE_USDT" default:"3"`
	PositionSize    float64 `envconfig:"POSITION_SIZE" default:"0.001"`

	Debug bool `envconfig:"POSITION_DEBUG" default:"true"`
}

func (c *config) Load() error {
	return envconfig.Process("", c)
}

var m *Module

func New(ctx context.Context) (*Module, error) {
	_, span := tracer.Start(ctx, "pkg.internal.position.New")
	defer span.End()

	if m != nil {
		return m, nil
	}

	l := logger.New("position")

	c := &config{}
	if err := c.Load(); err != nil {
		return nil, errors.Wrap(err, "loading configuration")
	}

	bm, err := bybit.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "creating bybit module")
	}

	cm, err := chart.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "creating chart module")
	}

	m = &Module{
		config: c,
		log:    l,

		chartModule: cm,

		bybitModule: bm,
	}

	return m, nil
}

type Module struct {
	config *config
	log    logger.Logger

	currentPosition *model.Position

	chartModule *chart.Module

	bybitModule *bybit.Module
}
