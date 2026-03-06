package usecases

import (
	"context"
	comercio "prueba-go/internal/domain/commerce"
	"prueba-go/pkg/types"
)

func (u *useCases) GetByID(ctx context.Context, id types.UID) (comercio.Comercio, error) {
	if id == "" {
		return comercio.Comercio{}, comercio.ErrInvalidComercioID
	}

	c, err := u.comercioRepo.GetByID(ctx, id)
	if err != nil {
		return comercio.Comercio{}, err
	}

	return c, nil
}
