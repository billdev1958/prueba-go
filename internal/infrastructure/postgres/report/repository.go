package report

import (
	"context"
	report "prueba-go/internal/domain/reports"
	"prueba-go/internal/infrastructure/postgres"
	"prueba-go/pkg/types"
	"prueba-go/pkg/util/money"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxReportRepository struct {
	postgres.PgxRepository
}

func NewReportRepository(pool *pgxpool.Pool) report.ReportRepository {
	return &pgxReportRepository{
		PgxRepository: postgres.PgxRepository{
			Pool: pool,
		},
	}
}

func (r *pgxReportRepository) GetGlobalEarnings(ctx context.Context) (money.Amount, error) {
	query := `
		SELECT COALESCE(SUM(net_amount), 0)
		FROM transactions
	`
	var amountStr string
	err := r.Pool.QueryRow(ctx, query).Scan(&amountStr)
	if err != nil {
		return money.Amount{}, err
	}

	return money.NewAmmount(amountStr)
}

func (r *pgxReportRepository) GetEarningsByMerchant(ctx context.Context, merchantID types.UID) (money.Amount, error) {
	query := `
		SELECT COALESCE(SUM(net_amount), 0)
		FROM transactions
		WHERE comercio_id = $1
	`
	var amountStr string
	err := r.Pool.QueryRow(ctx, query, merchantID).Scan(&amountStr)
	if err != nil {
		return money.Amount{}, err
	}

	return money.NewAmmount(amountStr)
}
