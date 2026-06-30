package handler

import (
	"encoding/json"
	"link-cutter/internal/app/util"
	"link-cutter/internal/link/service"
	"net/http"

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
		h.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestId(r)),
		)
		util.WriteError(
			w,
			http.StatusBadRequest,
			err.Error(),
			util.GetRequestId(r),
		)
		return
	}

	id, err := h.service.Create(ctx, createDto.Origin)
	if err != nil {
		h.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestId(r)),
		)
		util.WriteError(
			w,
			http.StatusInternalServerError,
			err.Error(),
			util.GetRequestId(r),
		)
		return
	}

	h.logger.Info("link created",
		zap.String("id", id.String()),
		zap.String("req_id", util.GetRequestId(r)),
	)
	w.WriteHeader(http.StatusCreated)
}

func (h *LinkHandler) Go(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("shortId")

	link, err := h.service.FindByShortId(ctx, id)
	if err != nil {
		h.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestId(r)),
		)
		util.WriteError(
			w,
			http.StatusInternalServerError,
			err.Error(),
			util.GetRequestId(r),
		)
		return
	}

	http.Redirect(w, r, link.Origin, http.StatusFound)

	h.logger.Info("successful link redirect",
		zap.String("shortId", id),
		zap.String("req_id", util.GetRequestId(r)),
	)
}
