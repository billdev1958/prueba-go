package usecases

import (
	"context"
	"prueba-go/internal/domain/audit"
	comercio "prueba-go/internal/domain/commerce"
	"prueba-go/pkg/types"
	"prueba-go/pkg/util/uuid"
	"time"
)

func (u *useCases) GetByID(ctx context.Context, id types.UID) (comercio.Comercio, error) {
	if id == "" {
		return comercio.Comercio{}, comercio.ErrInvalidComercioID
	}

	c, err := u.comercioRepo.GetByID(ctx, id)
	if err == nil {
		actor, _ := ctx.Value("actor").(string)
		if actor == "" {
			actor = "system_unknown"
		}

		go func(actor string, resourceID types.UID) {
			_ = u.auditRepo.Save(context.Background(), &audit.AuditLog{
				LogID:      uuid.NewUUID(),
				Action:     "COMMERCE_GET_BY_ID",
				Actor:      actor,
				ResourceID: resourceID,
				Timestamp:  time.Now(),
			})
		}(actor, id)
	}
	return c, err
}
