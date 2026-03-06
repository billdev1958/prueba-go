package usecases

import (
	"context"
	"prueba-go/pkg/util/money"
)

func (u *useCases) GetGlobalEarnings(ctx context.Context) (money.Amount, error) {
	return u.reportRepository.GetGlobalEarnings(ctx)
}
