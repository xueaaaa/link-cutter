package main

import (
	"context"
	"errors"
	"link-cutter/internal/app/config"
	"link-cutter/internal/link/handler"
	"link-cutter/internal/link/repository"
	"link-cutter/internal/link/service"
	"link-cutter/internal/postgres"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load("config.yaml")
	if err != nil {
		return
	}

	r := chi.NewRouter()

	dbConn, err := postgres.Connect(ctx, cfg.Postgres)
	if err != nil {
		return
	}

	linkRepo := repository.NewLinkRepository(dbConn)
	linkService := service.NewLinkService(linkRepo)
	linkHandler := handler.NewLinkHandler(linkService)

	r.Route("/link", func(r chi.Router) {
		r.Post("/", linkHandler.Create)
	})

	if err := http.ListenAndServe(":"+cfg.RunPort, r); !errors.Is(err, http.ErrServerClosed) {
	}
}
