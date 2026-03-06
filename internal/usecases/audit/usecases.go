package audit

import (
	"context"
	"prueba-go/internal/domain/audit"
)

type AuditUsecase interface {
	GetAll(ctx context.Context) ([]audit.AuditLog, error)
}

type useCases struct {
	repo audit.AuditRepository
}

func NewUsecases(repo audit.AuditRepository) AuditUsecase {
	return &useCases{
		repo: repo,
	}
}

func (u *useCases) GetAll(ctx context.Context) ([]audit.AuditLog, error) {
	return u.repo.GetAll(ctx)
}
