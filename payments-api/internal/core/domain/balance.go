package domain

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/core/constants"
	"github.com/jtonynet/go-payments-api/internal/core/customError"
	"github.com/jtonynet/go-payments-api/internal/core/port"
	"github.com/shopspring/decimal"
)

type Balance interface {
	ApproveTransaction(*Transaction) (Balance, *customError.CustomError)
	GetAccountID() uint
	GetAmountTotal() decimal.Decimal
	GetCategories() Categories
}

type BaseBalance struct {
	AccountID   uint
	AmountTotal decimal.Decimal
	Categories
}

type BalanceL1 struct {
	BaseBalance
}

func NewBalanceFromEntity(be port.BalanceEntity) (*BalanceL1, error) {
	categoryItens := make(map[int]Category)
	for _, ce := range be.Categories {
		category := Category{
			ID:       ce.ID,
			Name:     ce.Category.Name,
			Amount:   ce.Amount,
			MCCcodes: ce.Category.MCCcodes,
			Order:    ce.Category.Order,
		}
		categoryItens[ce.Category.Order] = category
	}

	if len(categoryItens) == 0 {
		return &BalanceL1{}, fmt.Errorf("dont instatiate with AccountID %v not found", be.AccountID)
	}

	categories := Categories{Itens: categoryItens}

	b := &BalanceL1{
		BaseBalance: BaseBalance{
			AccountID:   be.AccountID,
			AmountTotal: be.AmountTotal,
			Categories:  categories,
		},
	}

	return b, nil
}

func (b *BalanceL1) GetAccountID() uint {
	return b.AccountID
}

func (b *BalanceL1) GetAmountTotal() decimal.Decimal {
	return b.AmountTotal
}

func (b *BalanceL1) GetCategories() Categories {
	return b.Categories
}

func (b *BalanceL1) ApproveTransaction(tDomain *Transaction) (Balance, *customError.CustomError) {
	mcc := tDomain.MCCcode

	bCategoryToDebt, err := b.Categories.GetByMCC(mcc)
	if err != nil {
		return &BalanceL1{}, customError.New(constants.CODE_REJECTED_GENERIC, err.Error())
	}

	approvedBalance := bCategoryToDebt.Approve(tDomain)

	if !approvedBalance {
		msg := fmt.Sprintf("%s balance category has insuficient funds, %v available, transaction value is %v for account %v",
			bCategoryToDebt.Name,
			bCategoryToDebt.Amount,
			tDomain.TotalAmount,
			b.AccountID,
		)
		return &BalanceL1{}, customError.New(constants.CODE_REJECTED_INSUFICIENT_FUNDS, msg)
	}

	bCategoryToDebt.Amount = bCategoryToDebt.Amount.Sub(tDomain.TotalAmount)
	b.Categories.Itens[bCategoryToDebt.Order] = bCategoryToDebt

	b.AmountTotal = b.AmountTotal.Sub(tDomain.TotalAmount)

	return b, nil
}
