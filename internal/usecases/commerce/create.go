package usecases

import (
	"context"
	comercio "prueba-go/internal/domain/commerce"
	"prueba-go/pkg/util/money"
	"prueba-go/pkg/util/uuid"
)

func (u *useCases) Create(ctx context.Context, c *comercio.Comercio) (*comercio.Comercio, error) {
	c.ID = uuid.NewUUID()
	if c.Name == "" {
		return nil, comercio.ErrInvalidComercioName
	}

	maxRate, err := money.NewRate("1.00")
	if err != nil {
		return nil, err
	}

	if c.ComissionRate.IsNegative() || c.ComissionRate.IsGreaterThan(maxRate) {
		return nil, comercio.ErrInvalidComercioComissionRate
	}

	return u.comercioRepo.Create(ctx, c)
}
