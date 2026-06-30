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
	Create(ctx context.Context, origin string) (model.Link, error)
	FindByShortId(ctx context.Context, shortId string) (model.Link, error)
	Edit(ctx context.Context, link model.Link) error
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

func (s *linkService) Create(ctx context.Context, origin string) (model.Link, error) {
	link := model.Link{
		Origin:       origin,
		CreationDate: time.Now(),
	}

	for i := 0; i < 10; i++ {
		shortId := util.GenerateShortId(8)
		link.ShortId = shortId

		err := s.validate.StructCtx(ctx, link)
		if err != nil {
			return model.Link{}, err
		}

		id, err := s.repo.Create(ctx, repository.LinkModel(link))
		if err == nil {
			link.Id = id
			return link, nil
		}

		if errors.Is(err, errors2.ErrDuplicateShortId) {
			continue
		}

		return model.Link{}, err
	}

	return model.Link{}, errors2.ErrShortIdLimitExceeded
}

func (s *linkService) FindByShortId(ctx context.Context, shortId string) (model.Link, error) {
	lm, err := s.repo.FindByShortId(ctx, shortId)
	if err != nil {
		return model.Link{}, err
	}

	if lm == nil {
		return model.Link{}, errors2.ErrLinkNotFound
	}

	now := time.Now().UTC()
	lm.LastAccessDate = &now
	link := model.Link{
		Id:             lm.Id,
		ShortId:        lm.ShortId,
		Origin:         lm.Origin,
		CreationDate:   lm.CreationDate,
		LastAccessDate: lm.LastAccessDate,
	}

	err = s.Edit(ctx, link)
	if err != nil {
		return model.Link{}, err
	}

	return link, nil
}

func (s *linkService) Edit(ctx context.Context, link model.Link) error {
	lm := repository.LinkModel{
		Id:             link.Id,
		ShortId:        link.ShortId,
		Origin:         link.Origin,
		CreationDate:   link.CreationDate,
		LastAccessDate: link.LastAccessDate,
	}
	return s.repo.Edit(ctx, lm)
}
