package client

import (
	"context"
	"github.com/go-resty/resty/v2"
	"time"
)

type GrinexResponse struct {
	Timestamp int64      `json:"timestamp"`
	Asks      []RateItem `json:"asks"`
	Bids      []RateItem `json:"bids"`
}

type GrinexClient struct {
	client *resty.Client
}

func NewGrinexClient(url string) *GrinexClient {
	return &GrinexClient{
		client: resty.New().SetBaseURL(url).SetTimeout(5 * time.Second),
	}
}

func (c *GrinexClient) FetchRates(ctx context.Context) (*GrinexResponse, error) {
	var result GrinexResponse
	_, err := c.client.R().
		SetContext(ctx).
		SetResult(&result).
		Get("/api/v1/ticker") // Замените на реальный эндпоинт из документации
	return &result, err
}