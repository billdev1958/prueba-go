package transaction

import (
	"context"
	"prueba-go/internal/domain/audit"
	comercio "prueba-go/internal/domain/commerce"
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
	comercioRepo    comercio.ComercioRepository
	auditRepo       audit.AuditRepository
}

func NewUsecases(
	transactionRepo transaction.TransactionRepository,
	comercioRepo comercio.ComercioRepository,
	auditRepo audit.AuditRepository,
) TransactionUsecases {
	return &useCases{
		transactionRepo: transactionRepo,
		comercioRepo:    comercioRepo,
		auditRepo:       auditRepo,
	}
}
