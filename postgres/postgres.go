package postgres

import (
	"context"
	"fmt"
	"link-cutter/config"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context, cfg config.PostgresConfig) (*pgx.Conn, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
	)

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return conn, nil
}
