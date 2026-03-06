package transaction

import (
	"context"
	"prueba-go/internal/domain/audit"
	"prueba-go/internal/domain/transaction"
	"prueba-go/pkg/util/uuid"
	"time"
)

func (u *useCases) GetAll(ctx context.Context) ([]transaction.Transaction, error) {
	ts, err := u.transactionRepo.GetAll(ctx)
	if err == nil {
		actor, _ := ctx.Value("actor").(string)
		if actor == "" {
			actor = "system_unknown"
		}

		go func(actor string) {
			_ = u.auditRepo.Save(context.Background(), &audit.AuditLog{
				LogID:      uuid.NewUUID(),
				Action:     "TRANSACTION_GET_ALL",
				Actor:      actor,
				ResourceID: "ALL",
				Timestamp:  time.Now().UTC(),
			})
		}(actor)
	}
	return ts, err
}
