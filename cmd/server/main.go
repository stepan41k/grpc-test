package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/stepan41k/grpc-test/internal/app"
	"github.com/stepan41k/grpc-test/internal/config"
	"go.uber.org/zap"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting application")

	ctx := context.Background()

	application := app.New(ctx, log, cfg)

	go func() {
		if err := application.GRPCServer.Run(); err != nil {
			panic(err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	<-ctx.Done()

	log.Warn("stopping application")
	
	application.Close(ctx)

	log.Info("application stopped")
}

func setupLogger(env string) *zap.Logger {
	var log *zap.Logger
	var err error

	switch env {
	case envLocal:
		config := zap.NewDevelopmentConfig()
		config.DisableStacktrace = true
		log, err = config.Build()
		if err != nil {
			panic("failed to initialize local logger")
		}
	case envDev:
		log, err = zap.NewDevelopment()
		if err != nil {
			panic("failed to initialize development logger")
		}
	case envProd:
		log, err = zap.NewProduction()
		if err != nil {
			panic("failed to initialize production logger")
		}
	}

	return log
}
