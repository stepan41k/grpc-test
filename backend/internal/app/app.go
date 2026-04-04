package app

import (
	"context"

	grpcapp "github.com/stepan41k/grpc-test/internal/app/grpc"
	"github.com/stepan41k/grpc-test/internal/config"
	"github.com/stepan41k/grpc-test/internal/storage/postgres"

	"go.uber.org/zap"
)

type App struct {
	gRPCServer *grpcapp.App
	log        *zap.Logger
}

func New(ctx context.Context, log *zap.Logger, cfg *config.Config, connectionString string) *App {
	pool, err := postgres.New(ctx, connectionString)

	return &App{
		HTTPServer: httpApp,
		log:        log,
	}
}