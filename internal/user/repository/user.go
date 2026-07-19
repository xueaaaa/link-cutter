package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user UserModel) (pgtype.UUID, error)
	GetByEmail(ctx context.Context, email string) (UserModel, error)
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
