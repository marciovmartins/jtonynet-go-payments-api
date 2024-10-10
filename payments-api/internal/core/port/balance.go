package port

import "github.com/gofrs/uuid"

type BalanceDTORepository struct{}

type BalanceRepository interface {
	FindByAccountUUID(uuid.UUID) (BalanceDTORepository, error)
	Update(BalanceDTORepository) error
}
