package usecases

import (
	"context"
	"prueba-go/internal/domain/audit"
	report "prueba-go/internal/domain/reports"
	"prueba-go/pkg/types"
	"prueba-go/pkg/util/money"
)

type ReportUsecase interface {
	GetGlobalEarnings(ctx context.Context) (money.Amount, error)
	GetEarningsByMerchant(ctx context.Context, merchantID types.UID) (money.Amount, error)
}

type useCases struct {
	reportRepository report.ReportRepository
	auditRepo        audit.AuditRepository
}

func NewUsecases(reportRepository report.ReportRepository, auditRepo audit.AuditRepository) ReportUsecase {
	return &useCases{
		reportRepository: reportRepository,
		auditRepo:        auditRepo,
	}
}
