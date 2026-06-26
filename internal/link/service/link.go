package service

import (
	"context"
	"errors"
	errors2 "link-cutter/internal/app/errors"
	"link-cutter/internal/app/util"
	"link-cutter/internal/link/repository"
	"time"
)

type LinkService interface {
	Create(ctx context.Context, origin string) error
}

type linkService struct {
	repo repository.LinkRepository
}

func (s *linkService) Create(ctx context.Context, origin string) error {
	link := repository.LinkModel{
		Origin: origin,
		CreationDate: time.Now(),
	}

	for i := 0; i < 10; i++ {
		shortId := util.GenerateShortId(8)
		link.ShortId = shortId
		err := s.repo.Create(ctx, link)

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
