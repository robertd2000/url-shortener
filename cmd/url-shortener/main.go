package main

import (
	"fmt"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"

	"url-shortener/internal/http-server/handlers/url/save"
	mwLogger "url-shortener/internal/http-server/middleware/logger"

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
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}
	log.Info("server stopped")

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