package usecases

import (
	"context"
	"prueba-go/internal/domain/audit"
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
	auditRepo    audit.AuditRepository
}

func NewUsecases(comercioRepo comercio.ComercioRepository, auditRepo audit.AuditRepository) ComercioUsecase {
	return &useCases{
		comercioRepo: comercioRepo,
		auditRepo:    auditRepo,
	}
}
