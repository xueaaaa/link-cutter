package main

import (
	"context"
	"errors"
	"link-cutter/internal/app/config"
	logger2 "link-cutter/internal/app/logger"
	"link-cutter/internal/app/middleware"
	"link-cutter/internal/link/handler"
	"link-cutter/internal/link/repository"
	"link-cutter/internal/link/service"
	"link-cutter/internal/postgres"
	"net/http"

	"github.com/go-chi/chi/v5"
	middleware2 "github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx := context.Background()

	logger, cleanup, err := logger2.NewLogger()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	cfg, err := config.Load("config.yaml")
	if err != nil {
		logger.Error(err.Error())
		return
	}

	r := chi.NewRouter()

	dbConn, err := postgres.Connect(ctx, cfg.Postgres)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	linkRepo := repository.NewLinkRepository(dbConn)
	linkService := service.NewLinkService(linkRepo)
	linkHandler := handler.NewLinkHandler(linkService, logger)

	r.Use(middleware2.RequestID)
	r.Use(middleware.Logging(logger))

	r.Route("/link", func(r chi.Router) {
		r.Post("/", linkHandler.Create)
	})

	if err := http.ListenAndServe(":"+cfg.RunPort, r); !errors.Is(err, http.ErrServerClosed) {
		logger.Error(err.Error())
	}
}
