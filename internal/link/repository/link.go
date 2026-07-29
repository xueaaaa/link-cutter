package repository

import (
	"context"
	"errors"
	errors2 "link-cutter/internal/app/errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LinkRepository interface {
	Create(ctx context.Context, link LinkModel) (pgtype.UUID, error)
	FindById(ctx context.Context, id pgtype.UUID) (*LinkModel, error)
	FindByShortId(ctx context.Context, shortId string) (*LinkModel, error)
	Edit(ctx context.Context, link LinkModel) error
	EditLastAccess(ctx context.Context, linkId pgtype.UUID) error
	Delete(ctx context.Context, id pgtype.UUID) error
}

type linkRepository struct {
	db *pgxpool.Pool
}

func NewLinkRepository(db *pgxpool.Pool) LinkRepository {
	return &linkRepository{
		db: db,
	}
}

func (r *linkRepository) Create(ctx context.Context, link LinkModel) (pgtype.UUID, error) {
	sql := `INSERT INTO links (userId, shortId, origin, creationDate, lastAccessDate)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id;`

	var id pgtype.UUID
	err := r.db.QueryRow(
		ctx,
		sql,
		link.UserId,
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
		return pgtype.UUID{}, err
	} else if err != nil {
		return pgtype.UUID{}, err
	}

	return id, nil
}

func (r *linkRepository) findBy(ctx context.Context, fieldName string, key any) (*LinkModel, error) {
	sql := `SELECT id, userId, shortId, origin, creationDate, lastAccessDate FROM links
			WHERE ` + fieldName + `= $1`

	row := r.db.QueryRow(ctx, sql, key)
	var link LinkModel
	err := row.Scan(
		&link.Id,
		&link.UserId,
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

func (r *linkRepository) FindById(ctx context.Context, id pgtype.UUID) (*LinkModel, error) {
	return r.findBy(ctx, "id", id.String())
}

func (r *linkRepository) FindByShortId(ctx context.Context, shortId string) (*LinkModel, error) {
	return r.findBy(ctx, "shortId", shortId)
}

func (r *linkRepository) Edit(ctx context.Context, link LinkModel) error {
	sql := `UPDATE links
			SET origin = COALESCE(NULLIF($1, ''), origin)
			WHERE id = $2`

	tag, err := r.db.Exec(ctx, sql,
		link.Origin,
		link.Id,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors2.ErrLinkNotFound
	}
	return nil
}

func (r *linkRepository) EditLastAccess(ctx context.Context, id pgtype.UUID) error {
	sql := `UPDATE links
			SET lastAccessDate = $1
			WHERE id = $2`

	tag, err := r.db.Exec(ctx, sql,
		time.Now().UTC(),
		id,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors2.ErrLinkNotFound
	}
	return nil
}

func (r *linkRepository) Delete(ctx context.Context, id pgtype.UUID) error {
	sql := `DELETE FROM links WHERE id = $1`

	tag, err := r.db.Exec(ctx, sql, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors2.ErrLinkNotFound
	}
	return nil
}
