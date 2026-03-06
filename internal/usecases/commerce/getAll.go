package usecases

import (
	"context"
	comercio "prueba-go/internal/domain/commerce"
)

func (u *useCases) GetAll(ctx context.Context) ([]comercio.Comercio, error) {
	comercios, err := u.comercioRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(comercios) == 0 {
		return nil, err
	}

	return comercios, nil
}
