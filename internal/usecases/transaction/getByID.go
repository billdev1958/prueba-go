package transaction

import (
	"context"
	"prueba-go/internal/domain/transaction"
	"prueba-go/pkg/types"
)

func (u *useCases) GetByID(ctx context.Context, id types.UID) (transaction.Transaction, error) {
	if id == "" {
		return transaction.Transaction{}, transaction.ErrInvalidTransactionID
	}

	return u.transactionRepo.GetByID(ctx, id)
}
