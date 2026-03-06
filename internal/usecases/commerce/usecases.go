package usecases

import (
	"context"
	comercio "prueba-go/internal/domain/commerce"
	"prueba-go/pkg/types"
)

type ComercioUsecase interface {
	Create(ctx context.Context, c *comercio.Comercio) (*comercio.Comercio, error)
	Update(ctx context.Context, c *comercio.Comercio) error
	Delete(ctx context.Context, id types.UID) error
	GetByID(ctx context.Context, id types.UID) (comercio.Comercio, error)
	GetAll(ctx context.Context) ([]comercio.Comercio, error)
}

type useCases struct {
	comercioRepo comercio.ComercioRepository
}

func NewUsecases(comercioRepo comercio.ComercioRepository) ComercioUsecase {
	return &useCases{
		comercioRepo: comercioRepo,
	}
}
