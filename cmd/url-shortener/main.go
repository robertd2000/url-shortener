package main

import (
	"fmt"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

const (
	envLocal = "local"
 	envDev = "dev"
	envProd = "prod"
)

func main()  {
	cfg := config.MustLoad()
	
	fmt.Println(cfg)

	log := setupLogger(cfg.Env)

	log.Info("start", slog.String("env", cfg.Env))
	log.Debug("debug")

	storage, err := sqlite.New(cfg.Storage)

	fmt.Println(cfg.Storage)

	if err != nil {
		log.Error("failed to connect to db", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	_ = storage
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
		case envLocal:
			log = slog.New(
				slog.NewTextHandler(
					os.Stdout,
					&slog.HandlerOptions{Level: slog.LevelDebug}),
				)

		case envDev:
			log = slog.New(
				slog.NewJSONHandler(
					os.Stdout,
					 &slog.HandlerOptions{Level: slog.LevelDebug}),
				)

		case envProd:
			log = slog.New(
				slog.NewJSONHandler(
					os.Stdout,
						&slog.HandlerOptions{Level: slog.LevelInfo}),
				)

		default:
			log = slog.New(
				slog.NewJSONHandler(
					os.Stdout,
						&slog.HandlerOptions{Level: slog.LevelInfo}),
				)
	}

	return log
}