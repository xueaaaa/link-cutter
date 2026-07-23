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
	handler2 "link-cutter/internal/user/handler"
	repository2 "link-cutter/internal/user/repository"
	service2 "link-cutter/internal/user/service"
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

	userRepo := repository2.NewUserRepository(dbConn, logger)
	userService := service2.NewUserService(userRepo, cfg, logger)
	userHandler := handler2.NewUserHandler(userService, logger)

	r.Use(middleware2.RequestID)
	r.Use(middleware.Logging(logger))

	r.Get("/{shortId}", linkHandler.Go)
	r.Route("/link", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.OptionalAuth(cfg.JwtSigningKey, logger))
			r.Post("/", linkHandler.Create)
		})
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(cfg.JwtSigningKey, logger))
			r.Patch("/", linkHandler.Edit)
			r.Delete("/{shortId}", linkHandler.Delete)
		})
	})
	r.Route("/user", func(r chi.Router) {
		r.Post("/", userHandler.Create)
		r.Get("/", userHandler.Auth)
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(cfg.JwtSigningKey, logger))
			r.Patch("/", userHandler.Edit)
			r.Delete("/{id}", userHandler.Delete)
		})
	})

	if err := http.ListenAndServe(":"+cfg.RunPort, r); !errors.Is(err, http.ErrServerClosed) {
		logger.Error(err.Error())
	}
}
