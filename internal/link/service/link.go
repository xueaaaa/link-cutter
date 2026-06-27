package service

import (
	"context"
	"errors"
	errors2 "link-cutter/internal/app/errors"
	"link-cutter/internal/app/util"
	"link-cutter/internal/link/model"
	"link-cutter/internal/link/repository"
	"time"

	"github.com/go-playground/validator/v10"
)

type LinkService interface {
	Create(ctx context.Context, origin string) error
}

type linkService struct {
	repo     repository.LinkRepository
	validate *validator.Validate
}

func NewLinkService(repo repository.LinkRepository) LinkService {
	return &linkService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *linkService) Create(ctx context.Context, origin string) error {
	link := model.Link{
		Origin:       origin,
		CreationDate: time.Now(),
	}

	for i := 0; i < 10; i++ {
		shortId := util.GenerateShortId(8)
		link.ShortId = shortId

		err := s.validate.StructCtx(ctx, link)
		if err != nil {
			return err
		}

		err = s.repo.Create(ctx, repository.LinkModel(link))
		if err == nil {
			return nil
		}

		if errors.Is(err, errors2.ErrDuplicateShortId) {
			continue
		}

		return err
	}

	return errors2.ErrShortIdLimitExceeded
}
