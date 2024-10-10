package port

import "github.com/gofrs/uuid"

type TransactionDTORequest struct {
	AccountUUID uuid.UUID
}

type TransactionEntity struct {
}

type TransactionRepository interface {
	Save(TransactionEntity) error
}
