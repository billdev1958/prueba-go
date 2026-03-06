package transaction

import (
	"context"
	"prueba-go/internal/domain/transaction"
	"prueba-go/pkg/util/uuid"
)

func (u *useCases) Create(ctx context.Context, t *transaction.Transaction) (*transaction.Transaction, error) {
	t.ID = uuid.NewUUID()

	if t.Amount.IsZero() || t.Amount.IsNegative() {
		return nil, transaction.ErrInvalidTransactionAmount
	}

	commission, err := t.Amount.Mul(t.AppliedRate)
	if err != nil {
		return nil, err
	}

	t.Commission = commission

	netAmount, err := t.Amount.Sub(t.Commission)
	if err != nil {
		return nil, err
	}

	t.NetAmount = netAmount

	return u.transactionRepo.Create(ctx, t)

}
