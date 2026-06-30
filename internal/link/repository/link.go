package repository

import (
	"context"
	"errors"
	errors2 "link-cutter/internal/app/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type LinkRepository interface {
	Create(ctx context.Context, link LinkModel) (pgtype.UUID, error)
	FindByShortId(ctx context.Context, shortId string) (*LinkModel, error)
}

type linkRepository struct {
	db *pgx.Conn
}

func NewLinkRepository(db *pgx.Conn) LinkRepository {
	return &linkRepository{
		db: db,
	}
}

func (r *linkRepository) Create(ctx context.Context, link LinkModel) (pgtype.UUID, error) {
	sql := `INSERT INTO links (shortId, origin, creationDate, lastAccessDate)
			VALUES ($1, $2, $3, $4)
			RETURNING id;`

	var id pgtype.UUID
	err := r.db.QueryRow(
		ctx,
		sql,
		link.ShortId,
		link.Origin,
		link.CreationDate,
		link.LastAccessDate,
	).Scan(&id)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return pgtype.UUID{}, errors2.ErrDuplicateShortId
		}
	} else if err != nil {
		return pgtype.UUID{}, err
	}

	return id, nil
}

func (r *linkRepository) FindByShortId(ctx context.Context, shortId string) (*LinkModel, error) {
	sql := `SELECT (id, shortId, origin, creationDate, lastAccessAt) FROM links
			WHERE shortId = $1`

	row := r.db.QueryRow(ctx, sql, shortId)
	var link LinkModel
	err := row.Scan(
		&link.Id,
		&link.ShortId,
		&link.Origin,
		&link.CreationDate,
		&link.LastAccessDate)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &link, nil
}
