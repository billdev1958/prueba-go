package usecases

import (
	"context"
	"prueba-go/pkg/types"
	"prueba-go/pkg/util/money"
)

func (u *useCases) GetEarningsByMerchant(ctx context.Context, merchantID types.UID) (money.Amount, error) {
	return u.reportRepository.GetEarningsByMerchant(ctx, merchantID)
}
