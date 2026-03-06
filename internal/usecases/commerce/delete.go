package usecases

import (
	"context"
	comercio "prueba-go/internal/domain/commerce"
	"prueba-go/pkg/types"
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

	return u.comercioRepo.Delete(ctx, id)

}
