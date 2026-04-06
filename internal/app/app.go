package app

import (
	"context"
	"time"

	grpcapp "github.com/stepan41k/grpc-test/internal/app/grpc"
	"github.com/stepan41k/grpc-test/internal/client"
	"github.com/stepan41k/grpc-test/internal/config"
	"github.com/stepan41k/grpc-test/internal/metrics"
	"github.com/stepan41k/grpc-test/internal/service"
	"github.com/stepan41k/grpc-test/internal/storage/postgres"
	"github.com/stepan41k/grpc-test/internal/tracing"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"go.uber.org/zap"
)

type App struct {
	GRPCServer *grpcapp.App
	storage    *postgres.PGStorage
	tracerProvider *sdktrace.TracerProvider
	log        *zap.Logger
}

func New(ctx context.Context, log *zap.Logger, cfg *config.Config) *App {
	const path = "app.New"

	log = log.With(
		zap.String("path", path),
	)

	tp, err := tracing.InitTracer(ctx, cfg.OtelConfig.ServiceName, cfg.OtelConfig.URL)
	if err != nil {
		log.Fatal("failed to init tracer", zap.Error(err))
	}
	
	go func() {
		log.Info("starting prometheus metrics server", zap.String("addr", ":9090"))

		if err = metrics.StartMetricsServer(":9090"); err != nil {
			log.Error("prometheus metrics server failed", zap.Error(err))
		}
	}()

	connString2 := config.DTO(cfg)

	pool, err := postgres.New(ctx, connString2)
	if err != nil {
		log.Fatal("failed to init storage", zap.Error(err))
	}

	client := client.NewGrinexClient(cfg.GrinexConfig.URL)

	exchangeService := service.New(log, client, pool)

	grpcApp := grpcapp.New(log, exchangeService, cfg.ServerConfig.GRPCPort)

	return &App{
		GRPCServer: grpcApp,
		storage:    pool,
		tracerProvider: tp,
		log:        log,
	}
}

func (a *App) Close(ctx context.Context) {
	const path = "app.Stop"

	log := a.log.With(
		zap.String("op", path),
	)

	log.Warn("starting graceful shutdown")

	a.GRPCServer.Stop()
	
	if a.tracerProvider != nil {
        a.log.Warn("shutting down tracer")
        
        shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
        defer cancel()

        if err := a.tracerProvider.Shutdown(shutdownCtx); err != nil {
            a.log.Error("tracer shutdown failed", zap.Error(err))
        } else {
            a.log.Info("tracer stopped successfully")
        }
    }

	if a.storage != nil {
		log.Warn("closing database connection")
		a.storage.Close(ctx)
	}
	
	log.Info("database connection closing successfully")

	log.Info("graceful shutdown complete")
}
