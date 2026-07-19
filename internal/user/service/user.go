package service

import (
	"context"
	"link-cutter/internal/app/config"
	"link-cutter/internal/app/util"
	"link-cutter/internal/user/model"
	"link-cutter/internal/user/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(ctx context.Context, user model.User) (pgtype.UUID, error)
	Auth(ctx context.Context, email string, password string) (string, error)
}

type userSerivce struct {
	repo     repository.UserRepository
	config   config.Config
	validate *validator.Validate
}

func NewUserService(repo repository.UserRepository, config config.Config) UserService {
	return &userSerivce{
		repo:     repo,
		config:   config,
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

func (s *userSerivce) Auth(ctx context.Context, email, password string) (string, error) {
	err := s.validate.VarCtx(ctx, email, "email")
	if err != nil {
		return "", err
	}
	err = s.validate.VarCtx(ctx, password, "min=6,max=72,printascii,excludesall= ")
	if err != nil {
		return "", err
	}

	got, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	user := model.User(got)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	return util.IssueToken(s.config.JwtSigningKey, user)
}
