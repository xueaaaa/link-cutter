package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type LinkRepository interface {
	Create(ctx context.Context, link LinkModel) error
}

type linkRepository struct {
	db *pgx.Conn
}

func NewLinkRepository(db *pgx.Conn) LinkRepository {
	return &linkRepository{
		db: db,
	}
}

func (r *linkRepository) Create(ctx context.Context, link LinkModel) error {
	sql := `INSERT INTO links (shortId, origin, creationDate, lastAccessDate)
			VALUES ($1, $2, $3, $4);`

	_, err := r.db.Exec(
		ctx,
		sql,
		link.ShortId,
		link.Origin,
		link.CreationDate,
		link.LastAccessDate)

	return err
}
