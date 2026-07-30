package handler

import (
	"encoding/json"
	"link-cutter/internal/app/middleware"
	"link-cutter/internal/app/util"
	"link-cutter/internal/user/model"
	"link-cutter/internal/user/service"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
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
	ctx := util.PutRequestId(r.Context(), util.GetRequestId(r))
	var createDTO CreateUserDTO
	err := json.NewDecoder(r.Body).Decode(&createDTO)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}

	user := model.User{
		Email:    createDTO.Email,
		Username: createDTO.Username,
		Password: createDTO.Password,
	}
	id, err := h.service.Create(ctx, user)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	h.logger.Info("user created",
		zap.String("id", id.String()),
		zap.String("req_id", util.GetRequestId(r)),
	)
}

func (h *UserHandler) Auth(w http.ResponseWriter, r *http.Request) {
	ctx := util.PutRequestId(r.Context(), util.GetRequestId(r))
	var authDTO AuthUserDTO
	err := json.NewDecoder(r.Body).Decode(&authDTO)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}

	token, err := h.service.Auth(ctx, authDTO.Email, authDTO.Password)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(token))
	if err != nil {
		h.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestId(r)),
		)
	}
}

func (h *UserHandler) Edit(w http.ResponseWriter, r *http.Request) {
	ctx := util.PutRequestId(r.Context(), util.GetRequestId(r))
	var editDto EditUserDTO
	err := json.NewDecoder(r.Body).Decode(&editDto)

	if err != nil {
		util.HandleError(w, r, err)
		return
	}

	claims, _ := middleware.ClaimsFromContext(ctx)
	err = h.service.EnsureRights(ctx, claims.UserId, editDto.Id)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}

	user := model.User{
		Id:       editDto.Id,
		Username: editDto.Username,
		Password: editDto.Password,
	}
	err = h.service.Edit(ctx, user)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	h.logger.Info("successful user update",
		zap.String("id", user.Id.String()),
		zap.String("req_id", util.GetRequestId(r)),
	)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := util.PutRequestId(r.Context(), util.GetRequestId(r))
	rawId := r.PathValue("id")

	var id pgtype.UUID
	err := id.Scan(rawId)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}

	user, err := h.service.FindById(ctx, id)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}

	claims, _ := middleware.ClaimsFromContext(ctx)
	err = h.service.EnsureRights(ctx, claims.UserId, user.Id)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}

	err = h.service.Delete(ctx, user.Id)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	h.logger.Info("successful user delete",
		zap.String("id", user.Id.String()),
		zap.String("req_id", util.GetRequestId(r)),
	)
}
