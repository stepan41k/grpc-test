package app

import (
	"context"

	grpcapp "github.com/stepan41k/grpc-test/internal/app/grpc"
	"github.com/stepan41k/grpc-test/internal/client"
	"github.com/stepan41k/grpc-test/internal/config"
	"github.com/stepan41k/grpc-test/internal/service"
	"github.com/stepan41k/grpc-test/internal/storage/postgres"

	"go.uber.org/zap"
)

type App struct {
	GRPCServer *grpcapp.App
	log        *zap.Logger
}

func New(ctx context.Context, log *zap.Logger, cfg *config.Config) *App {
	connString2 := config.DTO(cfg)
	
	pool, err := postgres.New(ctx, connString2)
	if err != nil {
		panic(err)
	}

	client := client.NewGrinexClient("https://grinex.io")
	
	exchangeService := service.New(log, client, pool)

	grpcApp := grpcapp.New(log, exchangeService, cfg.ServerConfig.GRPCPort)

	return &App{
		GRPCServer: grpcApp,
		log:        log,
	}
}
