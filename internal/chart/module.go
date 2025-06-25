package chart

import (
	"context"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"github.com/kirill-a-belov/trader/internal/chart/model"
	"github.com/kirill-a-belov/trader/pkg/logger"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

type config struct {
	Port      string        `envconfig:"TRADER_CHART_PORT" default:"8080"`
	Depth     int           `envconfig:"TRADER_CHART_DEPTH" default:"60"`
	Timeframe time.Duration `envconfig:"TRADER_CHART_TIMEFRAME" default:"1m"`
}

func (c *config) Load() error {
	return envconfig.Process("", c)
}

var m *Module

func New(ctx context.Context) (*Module, error) {
	_, span := tracer.Start(ctx, "pkg.internal.chart.New")
	defer span.End()

	if m != nil {
		return m, nil
	}

	l := logger.New("chart")

	c := &config{}
	if err := c.Load(); err != nil {
		return nil, errors.Wrap(err, "loading configuration")
	}

	m = &Module{
		config: c,
		log:    l,

		candleList:    make([]*model.Candle, 0),
		markerStorage: make(map[time.Time][]*model.Marker),
	}

	return m, nil
}

type Module struct {
	config *config
	log    logger.Logger

	candleList    []*model.Candle
	markerStorage map[time.Time][]*model.Marker
}
