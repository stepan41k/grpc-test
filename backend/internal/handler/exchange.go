package handler

import (
	"context"

	"github.com/sudo-odner/minor/backend/services/auth_service/internal/service"
	pb "github.com/sudo-odner/minor/backend/services/auth_service/pkg/exchange_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Exchange interface {
	GetRates(ctx context.Context, req *pb.GetRatesRequest) (*pb.GetRatesResponse, error)
}

type ExchangeHandler struct {
    pb.UnimplementedExchangeServiceServer
    exchange Exchange
}

func Register(gRPCServer *grpc.Server, exchange Exchange) {
	pb.RegisterExchangeServiceServer(gRPCServer, &ExchangeHandler{exchange: exchange})
}

func (h *ExchangeHandler) GetRates(ctx context.Context, req *pb.GetRatesRequest) (*pb.GetRatesResponse, error) {
    // 1. Вызываем сервис
    res, err := h.service.ProcessRates(ctx, int(req.TopNIndex), ...)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }
    // 2. Формируем ответ
    return &pb.GetRatesResponse{...}, nil
}