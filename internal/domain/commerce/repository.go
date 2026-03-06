package comercio

import (
	"context"
	"prueba-go/pkg/types"
)

type ComercioRepository interface {
	Create(ctx context.Context, c *Comercio) (*Comercio, error)
	Update(ctx context.Context, c *Comercio) error
	Delete(ctx context.Context, id types.UID) error
	GetByID(ctx context.Context, id types.UID) (Comercio, error)
	GetByName(ctx context.Context, name string) (*Comercio, error)
	GetAll(ctx context.Context) ([]Comercio, error)
}
