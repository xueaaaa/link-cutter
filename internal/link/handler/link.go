package handler

import (
	"encoding/json"
	"errors"
	errors2 "link-cutter/internal/app/errors"
	"link-cutter/internal/app/util"
	"link-cutter/internal/link/model"
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

	created, err := h.service.Create(ctx, createDto.Origin)
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
		zap.String("id", created.Id.String()),
		zap.String("req_id", util.GetRequestId(r)),
	)

	data, err := json.Marshal(created)
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

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	if err != nil {
		h.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestId(r)),
		)
	}
}

func (h *LinkHandler) Go(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("shortId")

	link, err := h.service.FindByShortId(ctx, id)
	if err != nil {
		h.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestId(r)),
		)

		if errors.Is(err, errors2.ErrLinkNotFound) {
			util.WriteError(
				w,
				http.StatusNotFound,
				err.Error(),
				util.GetRequestId(r),
			)
		} else {
			util.WriteError(
				w,
				http.StatusInternalServerError,
				err.Error(),
				util.GetRequestId(r),
			)
		}
		return
	}

	http.Redirect(w, r, link.Origin, http.StatusFound)

	h.logger.Info("successful link redirect",
		zap.String("shortId", id),
		zap.String("req_id", util.GetRequestId(r)),
	)
}

func (h *LinkHandler) Edit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var editDto EditLinkDTO
	err := json.NewDecoder(r.Body).Decode(&editDto)

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

	link := model.Link{
		Id:     editDto.Id,
		Origin: editDto.Origin,
	}
	err = h.service.Edit(ctx, link)
	if err != nil {
		h.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestId(r)),
		)

		if errors.Is(err, errors2.ErrLinkNotFound) {
			util.WriteError(
				w,
				http.StatusNotFound,
				err.Error(),
				util.GetRequestId(r),
			)
		} else {
			util.WriteError(
				w,
				http.StatusInternalServerError,
				err.Error(),
				util.GetRequestId(r),
			)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	h.logger.Info("successful link update",
		zap.String("id", link.Id.String()),
		zap.String("req_id", util.GetRequestId(r)),
	)
}

func (h *LinkHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("shortId")

	link, err := h.service.FindByShortId(ctx, id)
	if err != nil {
		h.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestId(r)),
		)

		if errors.Is(err, errors2.ErrLinkNotFound) {
			util.WriteError(
				w,
				http.StatusNotFound,
				err.Error(),
				util.GetRequestId(r),
			)
		} else {
			util.WriteError(
				w,
				http.StatusInternalServerError,
				err.Error(),
				util.GetRequestId(r),
			)
		}
		return
	}

	err = h.service.Delete(ctx, link.Id)
	if err != nil {
		h.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestId(r)),
		)

		if errors.Is(err, errors2.ErrLinkNotFound) {
			util.WriteError(
				w,
				http.StatusNotFound,
				err.Error(),
				util.GetRequestId(r),
			)
		} else {
			util.WriteError(
				w,
				http.StatusInternalServerError,
				err.Error(),
				util.GetRequestId(r),
			)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
	h.logger.Info("successful link delete",
		zap.String("id", link.Id.String()),
		zap.String("req_id", util.GetRequestId(r)),
	)
}
