package domain

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/core/customError"
	"github.com/shopspring/decimal"
)

type Balance struct {
	AccountID   uint
	AmountTotal decimal.Decimal
	Categories
}

func (b *Balance) ApproveTransaction(tDomain *Transaction) (*Balance, *customError.CustomError) {
	mcc := tDomain.MCCcode

	bCategoryToDebt, err := b.Categories.GetByMCC(mcc)
	if err != nil {
		return &Balance{}, customError.New(CODE_REJECTED_GENERIC, err.Error())
	}

	approvedBalance := bCategoryToDebt.Approve(tDomain)
	if !approvedBalance {
		bCategoryFallback, err := b.getCategoryFallback()
		if err != nil {
			return &Balance{}, customError.New(CODE_REJECTED_GENERIC, err.Error())
		}

		bCategoryToDebt = bCategoryFallback
	}

	approvedBalanceFallback := bCategoryToDebt.Approve(tDomain)
	if !approvedBalanceFallback {
		msg := fmt.Sprintf("%s balance category has insuficient funds, %v available, transaction value is %v for account %v",
			bCategoryToDebt.Name,
			bCategoryToDebt.Amount,
			tDomain.TotalAmount,
			b.AccountID,
		)
		return &Balance{}, customError.New(CODE_REJECTED_INSUFICIENT_FUNDS, msg)
	}

	bCategoryToDebt.Amount = bCategoryToDebt.Amount.Sub(tDomain.TotalAmount)
	b.Categories.Itens[bCategoryToDebt.Order] = bCategoryToDebt

	b.AmountTotal = b.AmountTotal.Sub(tDomain.TotalAmount)

	return b, nil
}

func (b *Balance) getCategoryFallback() (Category, *customError.CustomError) {
	bCategoryToDebt, err := b.Categories.GetFallback()
	if err != nil {
		return Category{}, customError.New(CODE_REJECTED_GENERIC, err.Error())
	}
	return bCategoryToDebt, nil
}
