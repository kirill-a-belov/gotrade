package bybit

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"github.com/kirill-a-belov/trader/pkg/logger"
	"github.com/kirill-a-belov/trader/pkg/tracer"
)

type config struct {
	ApiRestUrl string `envconfig:"BYBIT_API_REST_URL" default:"https://api-testnet.bybit.com"`
	ApiWssUrl  string `envconfig:"BYBIT_API_WSS_URL" default:"wss://stream-testnet.bybit.com"`
	ApiKey     string `envconfig:"BYBIT_API_KEY"`
	ApiSecret  string `envconfig:"BYBIT_API_SECRET"`

	ReceiveWindow string `envconfig:"BYBIT_RECEIVE_WINDOW" default:"5000"`
}

func (c *config) Load() error {
	return envconfig.Process("", c)
}

var m *Module

func New(ctx context.Context) (*Module, error) {
	_, span := tracer.Start(ctx, "pkg.bybit.New")
	defer span.End()

	if m != nil {
		return m, nil
	}

	l := logger.New("bybit")

	c := &config{}
	if err := c.Load(); err != nil {
		return nil, errors.Wrap(err, "loading configuration")
	}

	m = &Module{
		config: c,
		log:    l,

		restClient: resty.New(),
		wssClient:  websocket.DefaultDialer,
	}

	return m, nil
}

type Module struct {
	config *config
	log    logger.Logger

	restClient *resty.Client
	wssClient  *websocket.Dialer
}
