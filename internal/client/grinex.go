package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stepan41k/grpc-test/internal/model"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type GrinexClient struct {
	http *resty.Client
}

func NewGrinexClient(url string) *GrinexClient {
	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	restyClient := resty.NewWithClient(httpClient).
		SetBaseURL(url).
		SetTimeout(5 * time.Second)

	return &GrinexClient{
		http: restyClient,
	}
}

func (c *GrinexClient) FetchRates(ctx context.Context) (*model.GrinexResponse, error) {
	var result model.GrinexResponse

	resp, err := c.http.R().
		SetContext(ctx).
		SetResult(&result).
		Get("/api/v1/spot/depth?symbol=usdta7a5")

	if err != nil {
		return nil, fmt.Errorf("request to grinex is faliled: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("api returned error status: %d", resp.StatusCode())
	}

	return &result, err
}
