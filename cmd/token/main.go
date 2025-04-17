package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	userHandler "github.com/stepan41k/MEDODS-Test/internal/http-server/handlers/user"
	authHandler "github.com/stepan41k/MEDODS-Test/internal/http-server/handlers/auth"
	authService "github.com/stepan41k/MEDODS-Test/internal/service/auth"
	userService "github.com/stepan41k/MEDODS-Test/internal/service/user"	
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/stepan41k/MEDODS-Test/cmd/migrator"
	"github.com/stepan41k/MEDODS-Test/internal/app"
	"github.com/stepan41k/MEDODS-Test/internal/config"
	"github.com/stepan41k/MEDODS-Test/internal/storage/postgres"
	"github.com/stepan41k/MEDODS-Test/internal/http-server/handlers/logout"
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

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	storagePath := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Username, cfg.Storage.DBName, os.Getenv("PG_DB_PASSWORD"), cfg.Storage.SSLMode)

	pool, err := postgres.New(context.Background(), storagePath)
	if err != nil {
		panic(err)
	}

	authService, userService := authService.New(log, pool), userService.New(log, pool)
	authHandler, userHandler := authHandler.New(log, authService), userHandler.New(log, userService)

	migrator.NewMigration("postgres://postgres:password12345@pgsql:5432/postgres?sslmode=disable", os.Getenv("MIGRATIONS_PATH"))

	router.Route("/user", func(r chi.Router) {
		r.Post("/new", userHandler.CreateUser(context.Background()))
	})

	router.Route("/auth", func(r chi.Router) {
		r.Post("/new/{guid}", authHandler.Create(context.Background()))
		r.Post("/refresh", authHandler.Refresh(context.Background()))
		r.Post("/logout", logout.Delete(context.Background(), log))
	})

	log.Info("starting server")

	application := app.New(log, cfg, router)

	go func() {
		application.HTTPServer.Run()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	signal := <-stop

	log.Info("stopping application", slog.String("signal", signal.String()))

	application.HTTPServer.Stop(context.Background())

	postgres.Close(context.Background(), pool)

	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {

	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}