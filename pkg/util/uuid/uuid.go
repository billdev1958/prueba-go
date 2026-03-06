package uuid

import (
	"github.com/google/uuid"
	"prueba-go/pkg/types"
)

func NewUUID() types.UID {
	return types.UID(uuid.New().String())
}
