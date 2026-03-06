package transaction

import (
	"context"
	"prueba-go/internal/domain/transaction"
	"prueba-go/pkg/types"
)

type TransactionUsecases interface {
	Create(ctx context.Context, t *transaction.Transaction) (*transaction.Transaction, error)
	GetByID(ctx context.Context, id types.UID) (transaction.Transaction, error)
	GetAll(ctx context.Context) ([]transaction.Transaction, error)
}

type useCases struct {
	transactionRepo transaction.TransactionRepository
}

func NewUsecases(transactionRepo transaction.TransactionRepository) TransactionUsecases {
	return &useCases{
		transactionRepo: transactionRepo,
	}
}
