package service

import (
	"fmt"
	"time"

	"github.com/jtonynet/go-payments-api/internal/core/domain"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

func mapTransactionRequestToMemoryLockEntity(tMemoryLock port.TransactionPaymentRequest) port.MemoryLockEntity {
	return port.MemoryLockEntity{
		Key:       tMemoryLock.AccountUID.String(),
		Timestamp: time.Now().UnixMilli(),
	}
}

func mapMerchantEntityToDomain(mEntity *port.MerchantEntity) domain.Merchant {
	return domain.Merchant{
		Name: mEntity.Name,
		MCC:  mEntity.MCC,
	}
}

func mapAccountEntityToDomain(aEntity port.AccountEntity) domain.Account {
	return domain.Account{
		ID:  aEntity.ID,
		UID: aEntity.UID,
	}
}

func mapBalanceEntityToDomain(bEntity port.BalanceEntity) (*domain.Balance, error) {
	categoryItens := make(map[int]domain.Category)
	for _, ce := range bEntity.Categories {
		category := domain.Category{
			ID:       ce.ID,
			Name:     ce.Category.Name,
			Amount:   ce.Amount,
			MCCs:     ce.Category.MCCs,
			Priority: ce.Category.Priority,
		}
		categoryItens[category.Priority] = category
	}

	if len(categoryItens) == 0 {
		return &domain.Balance{}, fmt.Errorf("failed to map Balance with AccountID %v, Cotegories not found", bEntity.AccountID)
	}

	categories := domain.Categories{Itens: categoryItens}

	b := &domain.Balance{
		AccountID:   bEntity.AccountID,
		AmountTotal: bEntity.AmountTotal,
		Categories:  categories,
	}

	return b, nil
}

func mapBalanceDomainToEntity(dBalance *domain.Balance) port.BalanceEntity {
	bCategories := make(map[int]port.BalanceByCategoryEntity)
	for _, categoryItem := range dBalance.Categories.Itens {
		bCategory := port.BalanceByCategoryEntity{
			ID:        categoryItem.ID,
			AccountID: dBalance.AccountID,
			Amount:    categoryItem.Amount,
			Category: port.CategoryEntity{
				Name:     categoryItem.Name,
				MCCs:     categoryItem.MCCs,
				Priority: categoryItem.Priority,
			},
		}

		bCategories[categoryItem.Priority] = bCategory
	}

	return port.BalanceEntity{
		AccountID:   dBalance.AccountID,
		AmountTotal: dBalance.AmountTotal,
		Categories:  bCategories,
	}
}

func mapTransactionDomainToEntity(tDomain *domain.Transaction) port.TransactionEntity {
	return port.TransactionEntity{
		AccountID:   tDomain.AccountID,
		TotalAmount: tDomain.TotalAmount,
		MCC:         tDomain.MCC,
		Merchant:    tDomain.Merchant,
	}
}
