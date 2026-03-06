package transaction

import (
	"context"
	"prueba-go/internal/domain/audit"
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
	auditRepo       audit.AuditRepository
}

func NewUsecases(transactionRepo transaction.TransactionRepository, auditRepo audit.AuditRepository) TransactionUsecases {
	return &useCases{
		transactionRepo: transactionRepo,
		auditRepo:       auditRepo,
	}
}
