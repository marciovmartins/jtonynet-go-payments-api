package port

import (
	"github.com/google/uuid"
)

type AccountEntity struct {
	ID  uint
	UID uuid.UUID
}

type AccountRepository interface {
	FindByUID(uuid.UUID) (AccountEntity, error)
}
