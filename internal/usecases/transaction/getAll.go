package transaction

import (
	"context"
	"prueba-go/internal/domain/transaction"
)

func (u *useCases) GetAll(ctx context.Context) ([]transaction.Transaction, error) {
	return u.transactionRepo.GetAll(ctx)
}
