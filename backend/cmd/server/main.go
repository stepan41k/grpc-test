package main

import (
	"os"
	"os/signal"
	"syscall"
	"fmt"

	"github.com/stepan41k/grpc-test/internal/app"
	"github.com/stepan41k/grpc-test/internal/config"
	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()
    
	log := setupLogger(cfg.Env)
    
	log.Info("starting application", zap.Any("config", cfg))
    
	storage := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.PostgreConfig.Host, cfg.PostgreConfig.Port, cfg.PostgreConfig.Username, cfg.PostgreConfig.DBName, os.Getenv("DB_PASSWORD"), cfg.PostgreConfig.SSLMode)
    
	application := app.New(log, cfg.GRPCPort, storage)
    
    
	go func() {
		application.GRPCServer.MustRun()
	}()
	
	//Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
    
	signal := <-stop
    
	log.Info("stopping application", slog.String("signal", signal.String()))
    
	application.GRPCServer.Stop()
    
	log.Info("application stopped")
    }
    

func setupLogger(env string) *zap.Logger {
	var log *zap.Logger
	var err error

	switch env {
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