package handler

import (
	"context"
	"time"

	"github.com/stepan41k/grpc-test/internal/grpc/pb"
	"github.com/stepan41k/grpc-test/internal/model"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ExchangeService interface {
	GetAndProcessRates(ctx context.Context, topIdx, n, m int) (*model.Result, error)
}

type ExchangeHandler struct {
	log *zap.Logger
	pb.UnimplementedExchangeServiceServer
	exchangeService ExchangeService
}

func Register(log *zap.Logger, gRPCServer *grpc.Server, exchangeService ExchangeService) {
	pb.RegisterExchangeServiceServer(gRPCServer, &ExchangeHandler{log: log, exchangeService: exchangeService})
}

func (eh *ExchangeHandler) GetRates(ctx context.Context, req *pb.GetRatesRequest) (*pb.GetRatesResponse, error) {
	const path = "handler.exchange.GetRates"
	
	log := eh.log.With(
		zap.String("path", path),
	)
	
	log.Info("start attempt getting and processing rates")
	
	data, err := eh.exchangeService.GetAndProcessRates(ctx, int(req.TopNIndex), int(req.AvgN), int(req.AvgM))
	if err != nil {
		log.Error("failed to get rates", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	
	log.Info("got and processed rates successfully")
	
	// If need a human-readable format, use time.Unix(data.Timestamp, 0).UTC().Format("2006-01-02 15:04:05")
	return &pb.GetRatesResponse{
		TopNPrice:  data.TopPrice,
		AvgNmPrice: data.AvgPrice,
		Timestamp: timestamppb.New(time.Unix(data.Timestamp, 0)),
	}, nil
}

func (eh *ExchangeHandler) Check(ctx context.Context, req *pb.CheckRequest) (*pb.CheckResponse, error) {
	return &pb.CheckResponse{Status: "SERVING"}, nil
}
