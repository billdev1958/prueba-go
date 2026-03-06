package usecases

import (
	"context"
	comercio "prueba-go/internal/domain/commerce"
	"prueba-go/pkg/util/money"
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

	return u.comercioRepo.Update(ctx, c)
}
