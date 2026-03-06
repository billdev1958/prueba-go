package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxRepository struct {
	Pool *pgxpool.Pool
}

func NewPgxRepository(pool *pgxpool.Pool) *PgxRepository {
	return &PgxRepository{
		Pool: pool,
	}
}
