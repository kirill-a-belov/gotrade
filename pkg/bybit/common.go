package bybit

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/kirill-a-belov/trader/pkg/tracer"
)

func (m *Module) signature(ctx context.Context, params string) string {
	_, span := tracer.Start(ctx, "pkg.bybit.signature")
	defer span.End()

	h := hmac.New(sha256.New, []byte(m.config.ApiSecret))
	h.Write([]byte(params))
	return hex.EncodeToString(h.Sum(nil))
}

func (m *Module) setHeaders(ctx context.Context, request *resty.Request, query string) {
	_, span := tracer.Start(ctx, "pkg.bybit.setHeaders")
	defer span.End()

	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	signatureBase := timestamp + m.config.ApiKey + m.config.ReceiveWindow + query
	signature := m.signature(ctx, signatureBase)

	request.SetHeader("X-BAPI-API-KEY", m.config.ApiKey)
	request.SetHeader("X-BAPI-SIGN", signature)
	request.SetHeader("X-BAPI-TIMESTAMP", timestamp)
	request.SetHeader("X-BAPI-RECV-WINDOW", m.config.ReceiveWindow)
}
