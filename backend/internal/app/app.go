package app

import (
	"github.com/go-chi/chi/v5"
	httpapp "github.com/stepan41k/grpc-test/internal/app/http"
	"github.com/stepan41k/grpc-test/internal/config"
	
	"go.uber.org/zap"
)

type App struct {
	HTTPServer *httpapp.App
	log        *zap.Logger
}

func New(log *zap.Logger, cfg *config.Config, router chi.Router) *App {
	httpApp := httpapp.New(log, cfg, router)
	app

	return &App{
		HTTPServer: httpApp,
		log:        log,
	}
}