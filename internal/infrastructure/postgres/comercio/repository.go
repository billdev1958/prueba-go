package comercio

import (
	"context"
	"errors"
	comercio "prueba-go/internal/domain/commerce"
	"prueba-go/internal/infrastructure/postgres"
	"prueba-go/pkg/types"
	"prueba-go/pkg/util/money"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxComercioRepository struct {
	postgres.PgxRepository
}

func NewComercioRepository(pool *pgxpool.Pool) comercio.ComercioRepository {
	return &pgxComercioRepository{
		PgxRepository: postgres.PgxRepository{
			Pool: pool,
		},
	}
}

func (r *pgxComercioRepository) Create(ctx context.Context, c *comercio.Comercio) (*comercio.Comercio, error) {
	query := `
		INSERT INTO comercios (id, name, comission_rate, created_at)
		VALUES ($1, $2, $3, $4)
	`
	c.CreatedAt = time.Now()
	_, err := r.Pool.Exec(ctx, query, c.ID, c.Name, c.ComissionRate.RateToString(), c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *pgxComercioRepository) Update(ctx context.Context, c *comercio.Comercio) error {
	query := `
		UPDATE comercios
		SET name = $1, comission_rate = $2
		WHERE id = $3
	`
	res, err := r.Pool.Exec(ctx, query, c.Name, c.ComissionRate.RateToString(), c.ID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return comercio.ErrComercioNotFound
	}
	return nil
}

func (r *pgxComercioRepository) Delete(ctx context.Context, id types.UID) error {
	query := `
		DELETE FROM comercios
		WHERE id = $1
	`
	res, err := r.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return comercio.ErrComercioNotFound
	}
	return nil
}

func (r *pgxComercioRepository) GetByID(ctx context.Context, id types.UID) (comercio.Comercio, error) {
	query := `
		SELECT id, name, comission_rate, created_at
		FROM comercios
		WHERE id = $1
	`
	var c comercio.Comercio
	var comissionRateStr string
	err := r.Pool.QueryRow(ctx, query, id).Scan(&c.ID, &c.Name, &comissionRateStr, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return comercio.Comercio{}, comercio.ErrComercioNotFound
		}
		return comercio.Comercio{}, err
	}

	rate, err := money.NewRate(comissionRateStr)
	if err != nil {
		return comercio.Comercio{}, err
	}
	c.ComissionRate = rate

	return c, nil
}

func (r *pgxComercioRepository) GetAll(ctx context.Context) ([]comercio.Comercio, error) {
	query := `
		SELECT id, name, comission_rate, created_at
		FROM comercios
	`
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comercios = []comercio.Comercio{}
	for rows.Next() {
		var c comercio.Comercio
		var comissionRateStr string
		if err := rows.Scan(&c.ID, &c.Name, &comissionRateStr, &c.CreatedAt); err != nil {
			return nil, err
		}

		rate, err := money.NewRate(comissionRateStr)
		if err != nil {
			return nil, err
		}
		c.ComissionRate = rate

		comercios = append(comercios, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comercios, nil
}
