package handler

import (
	"encoding/json"
	"link-cutter/internal/app/util"
	"link-cutter/internal/link/service"
	"net/http"
)

type LinkHandler struct {
	service service.LinkService
}

func NewLinkHandler(service service.LinkService) *LinkHandler {
	return &LinkHandler{
		service: service,
	}
}

func (h *LinkHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var createDto CreateLinkDTO
	err := json.NewDecoder(r.Body).Decode(&createDto)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Create(ctx, createDto.Origin)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}
