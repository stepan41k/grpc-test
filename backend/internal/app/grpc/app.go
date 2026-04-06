package http

import (
	"fmt"
	"net"

	"github.com/stepan41k/grpc-test/internal/handler"
	"github.com/stepan41k/grpc-test/internal/metrics"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	log        *zap.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *zap.Logger, exchangeService handler.ExchangeService, port int) *App {
	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			metrics.GrpcMetrics.UnaryServerInterceptor(),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	metrics.GrpcMetrics.InitializeMetrics(gRPCServer)

	handler.Register(log, gRPCServer, exchangeService)

	return &App{log: log, gRPCServer: gRPCServer, port: port}
}

func (a *App) Run() error {
	const path = "grpcapp.Run"

	log := a.log.With(
		zap.String("op", path),
		zap.Int("port", a.port),
	)

	go func() {
		a.log.Info("starting prometheus metrics server", zap.String("addr", ":9090"))

		if err := metrics.StartMetricsServer(":9090"); err != nil {
			a.log.Error("prometheus metrics server failed", zap.Error(err))
		}
	}()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	log.Info("starting grpc server", zap.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	return nil
}

func (a *App) Stop() {
	const path = "grpc.Stop"

	log := a.log.With(
		zap.String("op", path),
		zap.Int("port", a.port),
	)

	log.Warn("stopping grpc server")
	
	a.gRPCServer.GracefulStop()

	log.Info("grpc server stopped successfully")
}
