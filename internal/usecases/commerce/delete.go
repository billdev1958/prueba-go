package usecases

import (
	"context"
	"prueba-go/internal/domain/audit"
	comercio "prueba-go/internal/domain/commerce"
	"prueba-go/pkg/types"
	"prueba-go/pkg/util/uuid"
	"time"
)

func (u *useCases) Delete(ctx context.Context, id types.UID) error {
	if id == "" {
		return comercio.ErrComercioNotFound
	}

	_, err := u.comercioRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Posible Regla de Negocio: ¿Tiene comisiones asociadas?
	// TODO

	err = u.comercioRepo.Delete(ctx, id)
	if err == nil {
		actor, _ := ctx.Value("actor").(string)
		if actor == "" {
			actor = "system_unknown"
		}

		go func(actor string, resourceID types.UID) {
			_ = u.auditRepo.Save(context.Background(), &audit.AuditLog{
				LogID:      uuid.NewUUID(),
				Action:     "COMMERCE_DELETE",
				Actor:      actor,
				ResourceID: resourceID,
				Timestamp:  time.Now(),
			})
		}(actor, id)
	}
	return err
}
