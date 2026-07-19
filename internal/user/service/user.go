package service

import (
	"context"
	"link-cutter/internal/app/util"
	"link-cutter/internal/user/model"
	"link-cutter/internal/user/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserService interface {
	Create(ctx context.Context, user model.User) (pgtype.UUID, error)
}

type userSerivce struct {
	repo     repository.UserRepository
	validate *validator.Validate
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userSerivce{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *userSerivce) Create(ctx context.Context, user model.User) (pgtype.UUID, error) {
	user.CreationDate = time.Now()
	err := s.validate.StructCtx(ctx, user)
	if err != nil {
		return pgtype.UUID{}, err
	}

	hash, err := util.HashPassword(user.Password)
	if err != nil {
		return pgtype.UUID{}, err
	}
	user.Password = string(hash)

	return s.repo.Create(ctx, repository.UserModel(user))
}
