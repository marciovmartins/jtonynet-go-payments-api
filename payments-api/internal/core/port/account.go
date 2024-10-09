package port

import (
	"github.com/gofrs/uuid"
	gormModel "github.com/jtonynet/go-payments-api/internal/adapter/model/gormModel"
)

type AccountRepository interface {
	FindAll() ([]gormModel.Account, error)
	FindByUUID(uuid.UUID) (gormModel.Account, error)
}
