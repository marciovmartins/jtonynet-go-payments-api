package domain

import (
	"github.com/google/uuid"
)

type Account struct {
	ID  uint
	UID uuid.UUID
}
