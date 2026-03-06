package transaction

import (
	"context"
	"prueba-go/pkg/types"
)

type TransactionRepository interface {
	Create(ctx context.Context, t *Transaction) (*Transaction, error)
	Update(ctx context.Context, t *Transaction) (*Transaction, error)
	Delete(ctx context.Context, id types.UID) error
	GetByID(ctx context.Context, id types.UID) (Transaction, error)
	GetAll(ctx context.Context) ([]Transaction, error)
}
