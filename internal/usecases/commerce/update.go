package usecases

import (
	"context"
	"prueba-go/internal/domain/audit"
	comercio "prueba-go/internal/domain/commerce"
	"prueba-go/pkg/types"
	"prueba-go/pkg/util/money"
	"prueba-go/pkg/util/uuid"
	"time"
)

func (u *useCases) Update(ctx context.Context, c *comercio.Comercio) error {
	if c.ID == "" {
		return comercio.ErrComercioNotFound
	}

	if c.Name == "" {
		return comercio.ErrInvalidComercioName
	}
	maxRate, err := money.NewRate("1.00")
	if err != nil {
		return err
	}

	if c.ComissionRate.IsNegative() || c.ComissionRate.IsGreaterThan(maxRate) {
		return comercio.ErrInvalidComercioComissionRate
	}

	err = u.comercioRepo.Update(ctx, c)
	if err == nil {
		actor, _ := ctx.Value("actor").(string)
		if actor == "" {
			actor = "system_unknown"
		}

		go func(actor string, resourceID types.UID) {
			_ = u.auditRepo.Save(context.Background(), &audit.AuditLog{
				LogID:      uuid.NewUUID(),
				Action:     "COMMERCE_UPDATE",
				Actor:      actor,
				ResourceID: resourceID,
				Timestamp:  time.Now(),
			})
		}(actor, c.ID)
	}
	return err
}
