package usecases

import (
	"context"
	"errors"
	"prueba-go/internal/domain/audit"
	comercio "prueba-go/internal/domain/commerce"
	"prueba-go/pkg/types"
	"prueba-go/pkg/util/money"
	"prueba-go/pkg/util/uuid"
	"time"
)

func (u *useCases) Create(ctx context.Context, c *comercio.Comercio) (*comercio.Comercio, error) {
	c.ID = uuid.NewUUID()
	if c.Name == "" {
		return nil, comercio.ErrInvalidComercioName
	}

	existing, err := u.comercioRepo.GetByName(ctx, c.Name)
	if err == nil && existing != nil {
		return nil, comercio.ErrComercioAlreadyExists
	}
	if err != nil && !errors.Is(err, comercio.ErrComercioNotFound) {
		return nil, err
	}

	maxRate, err := money.NewRate("1.00")
	if err != nil {
		return nil, err
	}

	if c.ComissionRate.IsNegative() || c.ComissionRate.IsGreaterThan(maxRate) {
		return nil, comercio.ErrInvalidComercioComissionRate
	}

	res, err := u.comercioRepo.Create(ctx, c)
	if err == nil {
		actor := types.GetActor(ctx)

		go func(actor string, resourceID types.UID) {
			_ = u.auditRepo.Save(context.Background(), &audit.AuditLog{
				LogID:      uuid.NewUUID(),
				Action:     "COMMERCE_CREATE",
				Actor:      actor,
				ResourceID: resourceID,
				Timestamp:  time.Now().UTC(),
			})
		}(actor, res.ID)
	}
	return res, err
}
