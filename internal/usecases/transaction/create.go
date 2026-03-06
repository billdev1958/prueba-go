package transaction

import (
	"context"
	"prueba-go/internal/domain/audit"
	"prueba-go/internal/domain/transaction"
	"prueba-go/pkg/types"
	"prueba-go/pkg/util/uuid"
	"time"
)

func (u *useCases) Create(ctx context.Context, t *transaction.Transaction) (*transaction.Transaction, error) {
	t.ID = uuid.NewUUID()

	if t.Amount.IsZero() || t.Amount.IsNegative() {
		return nil, transaction.ErrInvalidTransactionAmount
	}

	comm, err := u.comercioRepo.GetByID(ctx, t.CommercioID)
	if err != nil {
		return nil, err
	}
	t.AppliedRate = comm.ComissionRate

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

	res, err := u.transactionRepo.Create(ctx, t)
	if err == nil {
		actor := types.GetActor(ctx)

		go func(actor string, resourceID types.UID) {
			_ = u.auditRepo.Save(context.Background(), &audit.AuditLog{
				LogID:      uuid.NewUUID(),
				Action:     "TRANSACTION_CREATE",
				Actor:      actor,
				ResourceID: resourceID,
				Timestamp:  time.Now().UTC(),
			})
		}(actor, res.ID)
	}

	return res, err
}
