package handler

import (
	"encoding/json"
	"link-cutter/internal/app/util"
	"link-cutter/internal/user/model"
	"link-cutter/internal/user/service"
	"net/http"

	"go.uber.org/zap"
)

type UserHandler struct {
	service service.UserService
	logger  *zap.Logger
}

func NewUserHandler(service service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var createDTO CreateUserDTO
	err := json.NewDecoder(r.Body).Decode(&createDTO)
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

	user := model.User{
		Email:    createDTO.Email,
		Username: createDTO.Username,
		Password: createDTO.Password,
	}
	id, err := h.service.Create(ctx, user)
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

	h.logger.Info("user created",
		zap.String("id", id.String()),
		zap.String("req_id", util.GetRequestId(r)),
	)
}

func (h *UserHandler) Auth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var authDTO AuthUserDTO
	err := json.NewDecoder(r.Body).Decode(&authDTO)
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

	token, err := h.service.Auth(ctx, authDTO.Email, authDTO.Password)
	if err != nil {
		h.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestId(r)),
		)
		util.WriteError(
			w,
			http.StatusUnauthorized,
			err.Error(),
			util.GetRequestId(r),
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(token))
	if err != nil {
		h.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestId(r)),
		)
		return
	}
}
