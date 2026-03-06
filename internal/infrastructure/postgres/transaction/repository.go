package transaction

import (
	"context"
	"errors"
	transaction "prueba-go/internal/domain/transaction"
	"prueba-go/internal/infrastructure/postgres"
	"prueba-go/pkg/types"
	"prueba-go/pkg/util/money"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxTransactionRepository struct {
	postgres.PgxRepository
}

func NewTransactionRepository(pool *pgxpool.Pool) transaction.TransactionRepository {
	return &pgxTransactionRepository{
		PgxRepository: postgres.PgxRepository{
			Pool: pool,
		},
	}
}

func (r *pgxTransactionRepository) Create(ctx context.Context, t *transaction.Transaction) (*transaction.Transaction, error) {
	query := `
		INSERT INTO transactions (id, comercio_id, amount, applied_rate, commission, net_amount, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	t.CreatedAt = time.Now()
	_, err := r.Pool.Exec(ctx, query,
		t.ID,
		t.CommercioID,
		t.Amount.AmmountToString(),
		t.AppliedRate.RateToString(),
		t.Commission.AmmountToString(),
		t.NetAmount.AmmountToString(),
		t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *pgxTransactionRepository) Update(ctx context.Context, t *transaction.Transaction) (*transaction.Transaction, error) {
	query := `
		UPDATE transactions
		SET comercio_id = $1, amount = $2, applied_rate = $3, commission = $4, net_amount = $5
		WHERE id = $6
	`
	res, err := r.Pool.Exec(ctx, query,
		t.CommercioID,
		t.Amount.AmmountToString(),
		t.AppliedRate.RateToString(),
		t.Commission.AmmountToString(),
		t.NetAmount.AmmountToString(),
		t.ID,
	)
	if err != nil {
		return nil, err
	}
	if res.RowsAffected() == 0 {
		return nil, transaction.ErrTransactionNotFound
	}
	return t, nil
}

func (r *pgxTransactionRepository) Delete(ctx context.Context, id types.UID) error {
	query := `
		DELETE FROM transactions
		WHERE id = $1
	`
	res, err := r.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return transaction.ErrTransactionNotFound
	}
	return nil
}

func (r *pgxTransactionRepository) GetByID(ctx context.Context, id types.UID) (transaction.Transaction, error) {
	query := `
		SELECT id, comercio_id, amount, applied_rate, commission, net_amount, created_at
		FROM transactions
		WHERE id = $1
	`
	var t transaction.Transaction
	var amountStr, appliedRateStr, commissionStr, netAmountStr string

	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&t.ID,
		&t.CommercioID,
		&amountStr,
		&appliedRateStr,
		&commissionStr,
		&netAmountStr,
		&t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return transaction.Transaction{}, transaction.ErrTransactionNotFound
		}
		return transaction.Transaction{}, err
	}

	amount, err := money.NewAmmount(amountStr)
	if err != nil {
		return transaction.Transaction{}, err
	}
	appliedRate, err := money.NewRate(appliedRateStr)
	if err != nil {
		return transaction.Transaction{}, err
	}
	commission, err := money.NewAmmount(commissionStr)
	if err != nil {
		return transaction.Transaction{}, err
	}
	netAmount, err := money.NewAmmount(netAmountStr)
	if err != nil {
		return transaction.Transaction{}, err
	}

	t.Amount = amount
	t.AppliedRate = appliedRate
	t.Commission = commission
	t.NetAmount = netAmount

	return t, nil
}

func (r *pgxTransactionRepository) GetAll(ctx context.Context) ([]transaction.Transaction, error) {
	query := `
		SELECT id, comercio_id, amount, applied_rate, commission, net_amount, created_at
		FROM transactions
	`
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions = []transaction.Transaction{}
	for rows.Next() {
		var t transaction.Transaction
		var amountStr, appliedRateStr, commissionStr, netAmountStr string
		err := rows.Scan(
			&t.ID,
			&t.CommercioID,
			&amountStr,
			&appliedRateStr,
			&commissionStr,
			&netAmountStr,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		amount, err := money.NewAmmount(amountStr)
		if err != nil {
			return nil, err
		}
		appliedRate, err := money.NewRate(appliedRateStr)
		if err != nil {
			return nil, err
		}
		commission, err := money.NewAmmount(commissionStr)
		if err != nil {
			return nil, err
		}
		netAmount, err := money.NewAmmount(netAmountStr)
		if err != nil {
			return nil, err
		}

		t.Amount = amount
		t.AppliedRate = appliedRate
		t.Commission = commission
		t.NetAmount = netAmount

		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
