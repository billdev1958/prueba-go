package comercio

import (
	"prueba-go/pkg/types"
	"prueba-go/pkg/util/money"
	"time"
)

type Comercio struct {
	ID            types.UID
	Name          string
	ComissionRate money.Rate
	CreatedAt     time.Time
}
