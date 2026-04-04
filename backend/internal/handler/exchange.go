package handler

type ExchangeHandler struct {
    pb.UnimplementedExchangeServiceServer
    service *service.ExchangeService // Ссылка на бизнес-логику
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