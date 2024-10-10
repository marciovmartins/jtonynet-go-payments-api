package port

import (
	"github.com/gofrs/uuid"
	"github.com/jtonynet/go-payments-api/internal/adapter/model/gormModel"
)

type AccountDTORepository struct {
	UUID uuid.UUID
}

type AccountRepository interface {
	FindAll() ([]gormModel.Account, error) //TODO: to remove, test purpose
	FindByUUID(uuid.UUID) (AccountDTORepository, error)
}
