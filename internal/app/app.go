package app

import (
	"log/slog"

	"github.com/go-chi/chi"
	httpapp "github.com/stepan41k/MEDODS-Test/internal/app/http"
	"github.com/stepan41k/MEDODS-Test/internal/config"
)

type App struct {
	log *slog.Logger
	HTTPServer *httpapp.App
}

func New(log *slog.Logger, cfg *config.Config, router chi.Router) *App {

	httpApp := httpapp.New(log, cfg, router)

	return &App{
		log: log,
		HTTPServer: httpApp,
	}
}