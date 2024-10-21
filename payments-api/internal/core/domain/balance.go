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
	next Balance
}

func (b *BaseBalance) GetAccountID() uint {
	return b.AccountID
}

func (b *BaseBalance) GetAmountTotal() decimal.Decimal {
	return b.AmountTotal
}

func (b *BaseBalance) GetCategories() Categories {
	return b.Categories
}

/*
	L1. Simple Authorizer:
		This  authorizer  approves  or  rejects a transaction based on the `MCC`.
		It maps the MCC from the `Transaction` parameter to a `Balance  Category
		to Debt`. If the category has sufficient funds, it debits the balance.
*/

type BalanceSimple struct {
	BaseBalance
}

/*
	L2. Authorizer with Fallback:
		This authorizer extends the `Simple Authorizer` one by adding a fallback mechanism.
		If the MCC category lacks funds  or  can't  be  mapped, it uses the  highest `Order`
		category without `MCC` codes  to  approve  the transaction, if sufficient funds are
		available.  Today,  this category  is `CASH`,  but  in  the  future,  it could   be
		`DREX`, `BITCOIN`, `BRICSCOIN (so surreal it's real)`, etc.
*/

type BalanceWithFallback struct {
	BaseBalance
}

func NewBalanceFromEntity(be port.BalanceEntity) (Balance, error) {
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
		return &BalanceSimple{}, fmt.Errorf("dont instatiate with AccountID %v not found", be.AccountID)
	}

	categories := Categories{Itens: categoryItens}

	baseBalance := BaseBalance{
		AccountID:   be.AccountID,
		AmountTotal: be.AmountTotal,
		Categories:  categories,
	}

	balanceSimple := &BalanceSimple{
		BaseBalance: baseBalance,
	}

	// Decoreted L2. Authorizer with Fallback
	baseBalance.next = balanceSimple
	balanceWithFallback := &BalanceWithFallback{
		BaseBalance: baseBalance,
	}

	return balanceWithFallback, nil
}

func (b *BalanceSimple) ApproveTransaction(tDomain *Transaction) (Balance, *customError.CustomError) {
	mcc := tDomain.MCCcode

	bCategoryToDebt, err := b.Categories.GetByMCC(mcc)
	if err != nil {
		return &BalanceSimple{}, customError.New(constants.CODE_REJECTED_GENERIC, err.Error())
	}

	approvedBalance := bCategoryToDebt.Approve(tDomain)
	if !approvedBalance {
		msg := fmt.Sprintf("%s balance category has insuficient funds, %v available, transaction value is %v for account %v",
			bCategoryToDebt.Name,
			bCategoryToDebt.Amount,
			tDomain.TotalAmount,
			b.AccountID,
		)
		return &BalanceSimple{}, customError.New(constants.CODE_REJECTED_INSUFICIENT_FUNDS, msg)
	}

	bCategoryToDebt.Amount = bCategoryToDebt.Amount.Sub(tDomain.TotalAmount)
	b.Categories.Itens[bCategoryToDebt.Order] = bCategoryToDebt

	b.AmountTotal = b.AmountTotal.Sub(tDomain.TotalAmount)

	return b, nil
}

func (b *BalanceWithFallback) ApproveTransaction(tDomain *Transaction) (Balance, *customError.CustomError) {
	balancePreviousStep, err := b.next.ApproveTransaction(tDomain)
	if err == nil {
		return balancePreviousStep, nil
	}

	bCategoryToDebt, cErr := b.Categories.GetFallback()
	if cErr != nil {
		msg := "balance category fallback not found"
		return &BalanceSimple{}, customError.New(constants.CODE_REJECTED_GENERIC, msg)
	}

	approvedBalance := bCategoryToDebt.Approve(tDomain)
	if !approvedBalance {
		msg := fmt.Sprintf("%s balance category has insuficient funds, %v available, transaction value is %v for account %v",
			bCategoryToDebt.Name,
			bCategoryToDebt.Amount,
			tDomain.TotalAmount,
			b.AccountID,
		)
		return &BalanceSimple{}, customError.New(constants.CODE_REJECTED_INSUFICIENT_FUNDS, msg)
	}

	bCategoryToDebt.Amount = bCategoryToDebt.Amount.Sub(tDomain.TotalAmount)
	b.Categories.Itens[bCategoryToDebt.Order] = bCategoryToDebt

	b.AmountTotal = b.AmountTotal.Sub(tDomain.TotalAmount)

	return b, nil
}
