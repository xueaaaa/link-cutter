package repository

import (
	"context"
	"errors"
	errors2 "link-cutter/internal/app/errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserRepository interface {
	Create(ctx context.Context, user UserModel) (pgtype.UUID, error)
	FindById(ctx context.Context, id pgtype.UUID) (*UserModel, error)
	FindByEmail(ctx context.Context, email string) (*UserModel, error)
	Edit(ctx context.Context, user UserModel) error
	EditLastAccess(ctx context.Context, userId pgtype.UUID) error
}

type userRepository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewUserRepository(db *pgxpool.Pool, logger *zap.Logger) UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

func (r *userRepository) Create(ctx context.Context, user UserModel) (pgtype.UUID, error) {
	sql := `INSERT INTO users (email, username, password, creationDate, lastAccessDate)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id;`

	var id pgtype.UUID
	err := r.db.QueryRow(
		ctx,
		sql,
		user.Email,
		user.Username,
		user.Password,
		user.CreationDate,
		user.LastAccessDate,
	).Scan(&id)

	if err != nil {
		return pgtype.UUID{}, err
	}
	return id, nil
}

func (r *userRepository) findBy(ctx context.Context, fieldName string, key any) (*UserModel, error) {
	sql := `SELECT id, email, username, password, creationDate, lastAccessDate FROM users
			WHERE ` + fieldName + ` = $1`

	var user UserModel
	row := r.db.QueryRow(ctx, sql, key)
	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreationDate,
		&user.LastAccessDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindById(ctx context.Context, id pgtype.UUID) (*UserModel, error) {
	return r.findBy(ctx, "id", id.String())
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*UserModel, error) {
	return r.findBy(ctx, "email", email)
}

func (r *userRepository) Edit(ctx context.Context, user UserModel) error {
	r.logger.Debug("", zap.Any("", user))

	sql := `UPDATE users
			SET username = COALESCE(NULLIF($1, ''), username),
				password = COALESCE(NULLIF($2, ''), password)
			WHERE id = $3`

	tag, err := r.db.Exec(ctx, sql,
		user.Username,
		user.Password,
		user.Id,
	)

	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors2.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) EditLastAccess(ctx context.Context, id pgtype.UUID) error {
	sql := `UPDATE users
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
		return errors2.ErrUserNotFound
	}
	return nil
}
