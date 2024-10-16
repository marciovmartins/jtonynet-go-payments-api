package port

import (
	"github.com/gofrs/uuid"
)

type AccountEntity struct {
	ID  int
	UID uuid.UUID
}

type AccountRepository interface {
	FindByUID(uuid.UUID) (AccountEntity, error)
}
