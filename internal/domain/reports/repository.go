package report

import (
	"context"
	"prueba-go/pkg/types"
	"prueba-go/pkg/util/money"
)

type ReportRepository interface {
	GetGlobalEarnings(ctx context.Context) (money.Amount, error)
	GetEarningsByMerchant(ctx context.Context, merchantID types.UID) (money.Amount, error)
}
