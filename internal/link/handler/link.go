package handler

import (
	"encoding/json"
	"link-cutter/internal/app/util"
	"link-cutter/internal/link/service"
	"net/http"

	middleware2 "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type LinkHandler struct {
	service service.LinkService
	logger  *zap.Logger
}

func NewLinkHandler(service service.LinkService, logger *zap.Logger) *LinkHandler {
	return &LinkHandler{
		service: service,
		logger:  logger,
	}
}

func (h *LinkHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var createDto CreateLinkDTO
	err := json.NewDecoder(r.Body).Decode(&createDto)
	if err != nil {
		h.logger.Error(err.Error())
		util.WriteError(
			w,
			http.StatusBadRequest,
			err.Error(),
			r.Context().Value(middleware2.RequestIDKey).(string),
		)
		return
	}

	id, err := h.service.Create(ctx, createDto.Origin)
	if err != nil {
		h.logger.Error(err.Error())
		util.WriteError(
			w,
			http.StatusInternalServerError,
			err.Error(),
			r.Context().Value(middleware2.RequestIDKey).(string),
		)
		return
	}

	h.logger.Info("link created",
		zap.String("id", id.String()),
		zap.String("req_id", r.Context().Value(middleware2.RequestIDKey).(string)),
	)
	w.WriteHeader(http.StatusCreated)
}
