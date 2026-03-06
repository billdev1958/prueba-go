package usecases

import (
	"context"
	"prueba-go/internal/domain/audit"
	comercio "prueba-go/internal/domain/commerce"
	"prueba-go/pkg/util/uuid"
	"time"
)

func (u *useCases) GetAll(ctx context.Context) ([]comercio.Comercio, error) {
	comercios, err := u.comercioRepo.GetAll(ctx)
	if err == nil {
		actor, _ := ctx.Value("actor").(string)
		if actor == "" {
			actor = "system_unknown"
		}

		go func(actor string) {
			_ = u.auditRepo.Save(context.Background(), &audit.AuditLog{
				LogID:      uuid.NewUUID(),
				Action:     "COMMERCE_GET_ALL",
				Actor:      actor,
				ResourceID: "ALL",
				Timestamp:  time.Now().UTC(),
			})
		}(actor)
	}

	if len(comercios) == 0 {
		return nil, err
	}
	return comercios, nil
}
