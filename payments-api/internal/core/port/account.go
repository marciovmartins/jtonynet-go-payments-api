package port

import (
	"github.com/google/uuid"
	"github.com/jtonynet/go-payments-api/internal/core/domain"
	"github.com/shopspring/decimal"
)

type AccountEntity struct {
	ID  uint
	UID uuid.UUID
}

type AccountRepository interface {
	FindByUID(uuid.UUID) (AccountEntity, error)
}

func (a *AccountEntity) NewTransaction(
	mccCode string,
	totalAmount decimal.Decimal,
	mDomain domain.Merchant,
) *domain.Transaction {
	correctMccCode := mccCode
	if mDomain.MappedMccCode != "" {
		correctMccCode = mDomain.MappedMccCode
	}

	return &domain.Transaction{
		AccountID:   a.ID,
		AccountUID:  a.UID,
		MccCode:     correctMccCode,
		TotalAmount: totalAmount,
		Merchant:    mDomain.Name,
	}
}
