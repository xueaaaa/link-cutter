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
	"github.com/jackc/pgx/v5/pgtype"
)

type LinkService interface {
	Create(ctx context.Context, origin string) (pgtype.UUID, error)
	FindByShortId(ctx context.Context, shortId string) (model.Link, error)
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

func (s *linkService) Create(ctx context.Context, origin string) (pgtype.UUID, error) {
	link := model.Link{
		Origin:       origin,
		CreationDate: time.Now(),
	}

	for i := 0; i < 10; i++ {
		shortId := util.GenerateShortId(8)
		link.ShortId = shortId

		err := s.validate.StructCtx(ctx, link)
		if err != nil {
			return pgtype.UUID{}, err
		}

		id, err := s.repo.Create(ctx, repository.LinkModel(link))
		if err == nil {
			return id, nil
		}

		if errors.Is(err, errors2.ErrDuplicateShortId) {
			continue
		}

		return pgtype.UUID{}, err
	}

	return pgtype.UUID{}, errors2.ErrShortIdLimitExceeded
}

func (s *linkService) FindByShortId(ctx context.Context, shortId string) (model.Link, error) {
	lm, err := s.repo.FindByShortId(ctx, shortId)
	if err != nil {
		return model.Link{}, err
	}

	if lm == nil {
		return model.Link{}, errors2.ErrLinkNotFound
	}

	link := model.Link{
		Id:             lm.Id,
		ShortId:        lm.ShortId,
		Origin:         lm.Origin,
		CreationDate:   lm.CreationDate,
		LastAccessDate: lm.LastAccessDate,
	}

	return link, nil
}
