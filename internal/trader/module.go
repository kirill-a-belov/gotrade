package trader

import (
	"context"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"github.com/kirill-a-belov/trader/internal/chart"
	"github.com/kirill-a-belov/trader/internal/simple_algo"
	"github.com/kirill-a-belov/trader/pkg/logger"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

type config struct{}

func (c *config) Load() error {
	return envconfig.Process("", c)
}

var m *Module

type tradingAlgorithm interface {
	Process(ctx context.Context) error
	Name() string
}

func New(ctx context.Context) (*Module, error) {
	_, span := tracer.Start(ctx, "pkg.internal.trader.New")
	defer span.End()

	if m != nil {
		return m, nil
	}

	l := logger.New("trader")

	c := &config{}
	if err := c.Load(); err != nil {
		return nil, errors.Wrap(err, "loading configuration")
	}

	sam, err := simple_algo.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "creating simple algo")
	}

	cm, err := chart.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "creating chart model")
	}

	m = &Module{
		config: c,
		log:    l,

		tradingAlgorithmList: []tradingAlgorithm{sam},
		chartModule:          cm,
	}

	return m, nil
}

type Module struct {
	config *config
	log    logger.Logger

	tradingAlgorithmList []tradingAlgorithm

	chartModule *chart.Module
}
