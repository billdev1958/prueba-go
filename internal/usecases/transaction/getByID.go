package transaction

import (
	"context"
	"prueba-go/internal/domain/audit"
	"prueba-go/internal/domain/transaction"
	"prueba-go/pkg/types"
	"prueba-go/pkg/util/uuid"
	"time"
)

func (u *useCases) GetByID(ctx context.Context, id types.UID) (transaction.Transaction, error) {
	if id == "" {
		return transaction.Transaction{}, transaction.ErrInvalidTransactionID
	}

	t, err := u.transactionRepo.GetByID(ctx, id)
	if err == nil {
		actor, _ := ctx.Value("actor").(string)
		if actor == "" {
			actor = "system_unknown"
		}

		go func(actor string, resourceID types.UID) {
			_ = u.auditRepo.Save(context.Background(), &audit.AuditLog{
				LogID:      uuid.NewUUID(),
				Action:     "TRANSACTION_GET_BY_ID",
				Actor:      actor,
				ResourceID: resourceID,
				Timestamp:  time.Now().UTC(),
			})
		}(actor, id)
	}
	return t, err
}
