package service

import (
	"context"
	"fmt"
	"time"

	"github.com/stepan41k/grpc-test/internal/lib/calculate"
	"github.com/stepan41k/grpc-test/internal/metrics"
	"github.com/stepan41k/grpc-test/internal/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("exchange-service")

type ExchangRepository interface {
	SaveRate(ctx context.Context, ask float64, bid float64, timestamp time.Time) error
}

type GrinexClient interface {
	FetchRates(ctx context.Context) (*model.GrinexResponse, error)
}

type ExchangeService struct {
	log                *zap.Logger
	client             GrinexClient
	exchangeRepository ExchangRepository
}

func New(log *zap.Logger, client GrinexClient, exchangeRepository ExchangRepository) *ExchangeService {
	return &ExchangeService{
		log:                log,
		client:             client,
		exchangeRepository: exchangeRepository,
	}
}

func (es *ExchangeService) GetAndProcessRates(ctx context.Context, topN, n, m int) (*model.Result, error) {
	const path = "service.exchange.GetAndProccessRates"

	log := es.log.With(
		zap.String("path", path),
	)

	ctx, span := tracer.Start(ctx, "GetAndProccessRates")
	defer span.End()

	span.SetAttributes(
		attribute.Int("request.top_n", topN),
		attribute.Int("request.avg_n", n),
		attribute.Int("request.avg_m", m),
	)

	log.Info("fetching data from Grinex API")

	data, err := es.client.FetchRates(ctx)
	if err != nil {
		metrics.ExternalAPIRequests.WithLabelValues("error").Inc()
		log.Error("failed to fetch data from Grinex:", zap.Error(err))
		return nil, err
	}
	
	log.Info("data from Grinex API fetched successfully")

	metrics.ExternalAPIRequests.WithLabelValues("success").Inc()

	span.SetAttributes(attribute.Int("exchange.items_received", len(data.Asks)))

	if len(data.Asks) == 0 || len(data.Bids) == 0 {
		log.Error("empty orderbook", zap.Error(err))
		return nil, fmt.Errorf("empty orderbook")
	}

	bestAsk := data.Asks[0].Price
	bestBid := data.Bids[0].Price

	metrics.LastUSDTPrice.Set(bestAsk)

	log.Info("attempting to save rate in database")

	err = es.exchangeRepository.SaveRate(ctx, bestAsk, bestBid, time.Unix(data.Timestamp, 0))
	if err != nil {
		log.Warn("failed to save rate into database:", zap.Error(err))
	}
	
	log.Info("attempting to calculate top rate")

	topPrice, err := calculate.CalculateTopN(data.Asks, topN)
	if err != nil {
		log.Error("failed to calculate topN:", zap.Error(err))
		return nil, err
	}
	
	log.Info("top rate calculated successfully")
	
	log.Info("attempting to calculate average rate")

	avgPrice, err := calculate.CalculateAvgNM(data.Asks, n, m)
	if err != nil {
		log.Warn("failed to calculate avgNM:", zap.Error(err))
		return nil, err
	}
	
	log.Info("average rate calculated successfully")

	return &model.Result{
		TopPrice:  topPrice,
		AvgPrice:  avgPrice,
		Timestamp: data.Timestamp,
	}, nil
}
