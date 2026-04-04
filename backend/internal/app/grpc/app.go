package http

import (
	"fmt"
	"net"

	"github.com/stepan41k/grpc-test/internal/handler"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	log        *zap.Logger
	gRPCServer *grpc.Server
	port int
}

func New(log *zap.Logger, exchangeService handler.Exchange, port int) *App {
	gRPCServer := grpc.NewServer()
	
	handler.Register(gRPCServer, exchangeService)

	return &App{log: log, gRPCServer: gRPCServer, port: port}
}

func (a *App) Run() error {
	const path = "grpcapp.Run"

	log := a.log.With(
		zap.String("op", path),
		zap.Int("port", a.port),
	)

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
	const path = "grpcapp.Stop"

	log := a.log.With(
		zap.String("op", path),
		zap.Int("port", a.port),
	)

	log.Info("stoping grpc server")

	a.gRPCServer.GracefulStop()

	log.Info("server stoped")
}