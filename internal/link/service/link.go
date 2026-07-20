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
	Create(ctx context.Context, link model.Link) (model.Link, error)
	FindById(ctx context.Context, id pgtype.UUID) (model.Link, error)
	FindByShortId(ctx context.Context, shortId string) (model.Link, error)
	Edit(ctx context.Context, link model.Link) error
	Delete(ctx context.Context, id pgtype.UUID) error
	EnsureRights(ctx context.Context, linkId pgtype.UUID, userId pgtype.UUID) error
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

func (s *linkService) Create(ctx context.Context, link model.Link) (model.Link, error) {
	link = model.Link{
		UserId:       link.UserId,
		Origin:       link.Origin,
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

func (s *linkService) findBy(ctx context.Context, find func() (*repository.LinkModel, error)) (model.Link, error) {
	lm, err := find()

	if err != nil {
		return model.Link{}, err
	}

	if lm == nil {
		return model.Link{}, errors2.ErrLinkNotFound
	}

	link := model.Link(*lm)
	err = s.Edit(ctx, link)
	if err != nil {
		return model.Link{}, err
	}

	return link, nil
}

func (s *linkService) FindById(ctx context.Context, id pgtype.UUID) (model.Link, error) {
	return s.findBy(ctx, func() (*repository.LinkModel, error) {
		return s.repo.FindById(ctx, id)
	})
}

func (s *linkService) FindByShortId(ctx context.Context, shortId string) (model.Link, error) {
	return s.findBy(ctx, func() (*repository.LinkModel, error) {
		return s.repo.FindByShortId(ctx, shortId)
	})
}

func (s *linkService) Edit(ctx context.Context, link model.Link) error {
	now := time.Now().UTC()
	link.LastAccessDate = &now

	lm := repository.LinkModel{
		Id:             link.Id,
		ShortId:        link.ShortId,
		Origin:         link.Origin,
		CreationDate:   link.CreationDate,
		LastAccessDate: link.LastAccessDate,
	}
	return s.repo.Edit(ctx, lm)
}

func (s *linkService) Delete(ctx context.Context, id pgtype.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *linkService) EnsureRights(ctx context.Context, linkId, userId pgtype.UUID) error {
	link, err := s.FindById(ctx, linkId)
	if err != nil {
		return err
	}
	if link.UserId != userId {
		return errors2.ErrNotEnoughRights
	}
	return nil
}
