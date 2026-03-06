package usecases

import (
	"context"
	"prueba-go/internal/domain/audit"
	"prueba-go/pkg/util/money"
	"prueba-go/pkg/util/uuid"
	"time"
)

func (u *useCases) GetGlobalEarnings(ctx context.Context) (money.Amount, error) {
	amount, err := u.reportRepository.GetGlobalEarnings(ctx)
	if err == nil {
		actor, _ := ctx.Value("actor").(string)
		if actor == "" {
			actor = "system_unknown"
		}

		go func(actor string) {
			_ = u.auditRepo.Save(context.Background(), &audit.AuditLog{
				LogID:      uuid.NewUUID(),
				Action:     "REPORT_GET_GLOBAL_EARNINGS",
				Actor:      actor,
				ResourceID: "GLOBAL",
				Timestamp:  time.Now(),
			})
		}(actor)
	}
	return amount, err
}
