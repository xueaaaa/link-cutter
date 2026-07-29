package service

import (
	"context"
	"link-cutter/internal/app/config"
	errors2 "link-cutter/internal/app/errors"
	"link-cutter/internal/app/util"
	"link-cutter/internal/user/model"
	"link-cutter/internal/user/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(ctx context.Context, user model.User) (pgtype.UUID, error)
	Auth(ctx context.Context, email string, password string) (string, error)
	FindById(ctx context.Context, id pgtype.UUID) (model.User, error)
	Edit(ctx context.Context, user model.User) error
	EnsureRights(ctx context.Context, ctxUserId pgtype.UUID, expectedUserId pgtype.UUID) error
	Delete(ctx context.Context, id pgtype.UUID) error
}

type userService struct {
	repo     repository.UserRepository
	config   config.Config
	logger   *zap.Logger
	validate *validator.Validate
}

func NewUserService(repo repository.UserRepository, config config.Config, logger *zap.Logger) UserService {
	return &userService{
		repo:     repo,
		config:   config,
		logger:   logger,
		validate: validator.New(),
	}
}

func (s *userService) Create(ctx context.Context, user model.User) (pgtype.UUID, error) {
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

func (s *userService) Auth(ctx context.Context, email, password string) (string, error) {
	err := s.validate.VarCtx(ctx, email, "email")
	if err != nil {
		return "", err
	}
	err = s.validate.VarCtx(ctx, password, "min=6,max=72,printascii,excludesall= ")
	if err != nil {
		return "", err
	}

	got, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if got == nil {
		return "", errors2.ErrUserNotFound
	}

	user := model.User(*got)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	t := time.Now()
	got.LastAccessDate = &t
	err = s.repo.EditLastAccess(ctx, got.Id)
	if err != nil {
		return "", err
	}

	return util.IssueToken(s.config.JwtSigningKey, user)
}

func (s *userService) FindById(ctx context.Context, id pgtype.UUID) (model.User, error) {
	um, err := s.repo.FindById(ctx, id)

	if err != nil {
		return model.User{}, err
	}

	if um == nil {
		return model.User{}, errors2.ErrUserNotFound
	}

	user := model.User(*um)
	err = s.repo.EditLastAccess(ctx, user.Id)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *userService) Edit(ctx context.Context, user model.User) error {
	if user.Username != "" {
		err := s.validate.VarCtx(ctx, user.Username, "min=4,max=16")
		if err != nil {
			return err
		}
	}
	if user.Password != "" {
		err := s.validate.VarCtx(ctx, user.Password, "min=6,max=72,printascii,excludesall= ")
		if err != nil {
			return err
		}

		p, err := util.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = string(p)
	}

	um := repository.UserModel{
		Id:       user.Id,
		Username: user.Username,
		Password: user.Password,
	}

	return s.repo.Edit(ctx, um)
}

func (s *userService) EnsureRights(ctx context.Context, ctxUserId pgtype.UUID, expectedUserId pgtype.UUID) error {
	user, err := s.FindById(ctx, expectedUserId)
	if err != nil {
		return err
	}
	if user.Id != ctxUserId {
		return errors2.ErrNotEnoughRights
	}
	return nil
}

func (s *userService) Delete(ctx context.Context, id pgtype.UUID) error {
	return s.repo.Delete(ctx, id)
}
