package transaction

import (
	"prueba-go/pkg/types"
	"prueba-go/pkg/util/money"
	"time"
)

type Transaction struct {
	ID          types.UID
	CommercioID types.UID
	Amount      money.Amount
	AppliedRate money.Rate
	Commission  money.Amount
	NetAmount   money.Amount
	CreatedAt   time.Time
}
