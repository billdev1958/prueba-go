package usecases

import (
	"context"
	"prueba-go/internal/domain/audit"
	"prueba-go/pkg/types"
	"prueba-go/pkg/util/money"
	"prueba-go/pkg/util/uuid"
	"time"
)

func (u *useCases) GetEarningsByMerchant(ctx context.Context, merchantID types.UID) (money.Amount, error) {
	amount, err := u.reportRepository.GetEarningsByMerchant(ctx, merchantID)
	if err == nil {
		actor, _ := ctx.Value("actor").(string)
		if actor == "" {
			actor = "system_unknown"
		}

		go func(actor string, resourceID types.UID) {
			_ = u.auditRepo.Save(context.Background(), &audit.AuditLog{
				LogID:      uuid.NewUUID(),
				Action:     "REPORT_GET_EARNINGS_BY_MERCHANT",
				Actor:      actor,
				ResourceID: resourceID,
				Timestamp:  time.Now(),
			})
		}(actor, merchantID)
	}
	return amount, err
}
