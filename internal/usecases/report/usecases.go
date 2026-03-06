package usecases

import (
	"context"
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
}

func NewUsecases(reportRepository report.ReportRepository) ReportUsecase {
	return &useCases{
		reportRepository: reportRepository,
	}
}
