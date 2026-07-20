package repository

import (
	"context"
	"link-cutter/internal/app/errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user UserModel) (pgtype.UUID, error)
	GetByEmail(ctx context.Context, email string) (UserModel, error)
	Edit(ctx context.Context, user UserModel) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db: db,
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

func (r *userRepository) GetByEmail(ctx context.Context, email string) (UserModel, error) {
	sql := `SELECT id, email, username, password, creationDate, lastAccessDate FROM users
			WHERE email = $1`

	var user UserModel
	row := r.db.QueryRow(ctx, sql, email)
	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreationDate,
		&user.LastAccessDate,
	)
	if err != nil {
		return UserModel{}, err
	}
	return user, nil
}

func (r *userRepository) Edit(ctx context.Context, user UserModel) error {
	sql := `UPDATE users
			SET email = $1, username = $2, password = $3, lastAccessDate = $4
			WHERE id = $5`

	tag, err := r.db.Exec(ctx, sql,
		user.Email,
		user.Username,
		user.Password,
		user.LastAccessDate,
		user.Id,
	)

	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.ErrUserNotFound
	}
	return nil
}
