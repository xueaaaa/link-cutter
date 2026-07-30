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
		s.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestIdFromContext(ctx)))
		return pgtype.UUID{}, errors2.ErrInvalidInputData
	}

	hash, err := util.HashPassword(user.Password)
	if err != nil {
		s.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestIdFromContext(ctx)))
		return pgtype.UUID{}, errors2.ErrFailedHashPassword
	}
	user.Password = string(hash)

	id, err := s.repo.Create(ctx, repository.UserModel(user))
	if err != nil {
		s.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestIdFromContext(ctx)))
		return pgtype.UUID{}, errors2.ErrDatabase
	}
	return id, nil
}

func (s *userService) Auth(ctx context.Context, email, password string) (string, error) {
	err := s.validate.VarCtx(ctx, email, "email")
	if err != nil {
		s.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestIdFromContext(ctx)))
		return "", errors2.ErrIncorrectAuthData
	}
	err = s.validate.VarCtx(ctx, password, "min=6,max=72,printascii,excludesall= ")
	if err != nil {
		s.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestIdFromContext(ctx)))
		return "", errors2.ErrIncorrectAuthData
	}

	got, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		s.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestIdFromContext(ctx)))
		return "", errors2.ErrDatabase
	}
	if got == nil {
		return "", errors2.ErrIncorrectAuthData
	}

	user := model.User(*got)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		s.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestIdFromContext(ctx)))
		return "", errors2.ErrIncorrectAuthData
	}

	t := time.Now()
	got.LastAccessDate = &t
	err = s.repo.EditLastAccess(ctx, got.Id)
	if err != nil {
		s.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestIdFromContext(ctx)))
		return "", errors2.ErrDatabase
	}

	token, err := util.IssueToken(s.config.JwtSigningKey, user)
	if err != nil {
		s.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestIdFromContext(ctx)))
		return "", errors2.ErrIssueToken
	}
	return token, nil
}

func (s *userService) FindById(ctx context.Context, id pgtype.UUID) (model.User, error) {
	um, err := s.repo.FindById(ctx, id)

	if err != nil {
		s.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestIdFromContext(ctx)))
		return model.User{}, errors2.ErrDatabase
	}

	if um == nil {
		return model.User{}, errors2.ErrUserNotFound
	}

	user := model.User(*um)
	err = s.repo.EditLastAccess(ctx, user.Id)
	if err != nil {
		return model.User{}, errors2.ErrDatabase
	}

	return user, nil
}

func (s *userService) Edit(ctx context.Context, user model.User) error {
	if user.Username != "" {
		err := s.validate.VarCtx(ctx, user.Username, "min=4,max=16")
		if err != nil {
			return errors2.ErrInvalidInputData
		}
	}
	if user.Password != "" {
		err := s.validate.VarCtx(ctx, user.Password, "min=6,max=72,printascii,excludesall= ")
		if err != nil {
			return errors2.ErrInvalidInputData
		}

		p, err := util.HashPassword(user.Password)
		if err != nil {
			s.logger.Error(err.Error(),
				zap.String("req_id", util.GetRequestIdFromContext(ctx)))
			return errors2.ErrFailedHashPassword
		}
		user.Password = string(p)
	}

	um := repository.UserModel{
		Id:       user.Id,
		Username: user.Username,
		Password: user.Password,
	}

	err := s.repo.Edit(ctx, um)
	if err != nil {
		s.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestIdFromContext(ctx)))
		return errors2.ErrDatabase
	}
	return nil
}

func (s *userService) EnsureRights(ctx context.Context, ctxUserId pgtype.UUID, expectedUserId pgtype.UUID) error {
	user, err := s.FindById(ctx, expectedUserId)
	if err != nil {
		s.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestIdFromContext(ctx)))
		return err
	}
	if user.Id != ctxUserId {
		return errors2.ErrNotEnoughRights
	}
	return nil
}

func (s *userService) Delete(ctx context.Context, id pgtype.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error(err.Error(),
			zap.String("req_id", util.GetRequestIdFromContext(ctx)))
		return errors2.ErrDatabase
	}
	return nil
}
