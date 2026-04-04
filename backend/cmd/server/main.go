package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
    // импорты репозитория, конфига и OTel
)

func main() {
	
	// 1. Load Config (Flags + Env)
    cfg := config.Load()

    // 2. Init OpenTelemetry
    tp := initTracer(cfg.OTEL_URL)
    defer tp.Shutdown(context.Background())

    // 3. Init DB & Migrations
    db := initDB(cfg.DatabaseURL)
    defer db.Close()

    // 4. Setup gRPC
    lis, err := net.Listen("tcp", cfg.Port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer(
        // Добавление OTel интерцептора
    )
    
    // Register Healthcheck
    healthCheck := health.NewServer()
    healthpb.RegisterHealthServer(s, healthCheck)
    
    // Register Business Service
    // ...

    // 5. Graceful Shutdown
    go func() {
        if err := s.Serve(lis); err != nil {
            log.Fatalf("failed to serve: %v", err)
        }
    }()

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
    <-stop

    log.Println("Shutting down gracefully...")
    s.GracefulStop()
	
	cfg := config.MustLoad()
	log := setupLogger(envDev)

	storagePath := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.PostgreConfig.Host, cfg.PostgreConfig.Port, cfg.PostgreConfig.Username, cfg.PostgreConfig.DBName, os.Getenv("POSTGRES_PASSWORD"), cfg.PostgreConfig.SSLMode)

	dbConn, err := postgres.New(context.Background(), storagePath)
	if err != nil {
		panic("failed to initialize DB connection")
	}

	log.Info("starting authentication service")

	authService := authService.New(dbConn, log)
	authHandler := authHandler.New(authService, log)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(cors.NewCORS)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/auth-service", func(r chi.Router) {
		r.Post("/register", authHandler.Register(context.Background()))
		r.Post("/login", authHandler.Login(context.Background()))
	})

	// router.Route("/token", func(r chi.Router) {
	// 	r.Post("/refresh")
	// })

	application := app.New(log, cfg, router)

	go func() {
		application.HTTPServer.Run()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	signal := <-stop

	log.Info("stopping application", zap.String("signal", signal.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	application.HTTPServer.Stop(ctx)

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