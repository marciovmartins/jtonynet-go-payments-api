package domain

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/support"
	"github.com/shopspring/decimal"
)

type Balance struct {
	AccountID   uint
	AmountTotal decimal.Decimal
	Categories

	Logger support.Logger
}

func (b *Balance) ApproveTransaction(tDomain *Transaction) (*Balance, *CustomError) {
	debitedBalanceCategories, err := b.getDebitedBalanceCategories(tDomain)
	if err != nil {
		return b, err
	}

	for priority, balanceCategory := range debitedBalanceCategories {
		b.Categories.Itens[priority] = balanceCategory
	}

	return b, nil
}

func (b *Balance) getDebitedBalanceCategories(tDomain *Transaction) (map[int]Category, *CustomError) {
	amountDebtRemaining := tDomain.TotalAmount
	debitedBalanceCategories := make(map[int]Category)

	bCategoryMCC, err := b.Categories.GetByMCC(tDomain.MCC)
	if err != nil {
		b.debugLog(
			fmt.Sprintf(
				"Error retrieving category by MCC, attempting to debit remaining amount (%s) directly from fallback.",
				amountDebtRemaining,
			),
		)

	} else if bCategoryMCC.Amount.GreaterThanOrEqual(tDomain.TotalAmount) {
		b.debugLog(
			fmt.Sprintf(
				"Sufficient funds available in category '%s' to cover the transaction amount.",
				bCategoryMCC.Name,
			),
		)

		amountDebtRemaining = decimal.NewFromFloat(0)

		bCategoryMCC.Amount = bCategoryMCC.Amount.Sub(tDomain.TotalAmount)
		debitedBalanceCategories[bCategoryMCC.Priority] = bCategoryMCC

	} else if bCategoryMCC.Amount.IsPositive() {
		b.debugLog(
			fmt.Sprintf(
				"Category '%s' has a positive balance but insufficient funds for the transaction. Check fallback category.",
				bCategoryMCC.Name,
			),
		)

		amountDebtRemaining = amountDebtRemaining.Sub(bCategoryMCC.Amount)

		bCategoryMCC.Amount = decimal.NewFromFloat(0)
		debitedBalanceCategories[bCategoryMCC.Priority] = bCategoryMCC
	}

	if amountDebtRemaining.GreaterThan(decimal.Zero) {
		bCategoryFallback, err := b.Categories.GetFallback()
		if err != nil {
			b.debugLog(
				"Error retrieving fallback category. Rolling back changes and returning an error.",
			)

			return make(map[int]Category), NewCustomError(CODE_REJECTED_GENERIC, "Category Fallback not found")

		} else if bCategoryFallback.Amount.GreaterThanOrEqual(amountDebtRemaining) {
			b.debugLog(
				"Fallback category has sufficient funds available for the transaction.",
			)

			bCategoryFallback.Amount = bCategoryFallback.Amount.Sub(amountDebtRemaining)
			debitedBalanceCategories[bCategoryFallback.Priority] = bCategoryFallback

			amountDebtRemaining = decimal.NewFromFloat(0)
		} else {
			b.debugLog(
				"Insufficient funds in the fallback category. Rolling back changes and returning an error.",
			)

			debitedBalanceCategories = make(map[int]Category)
			amountDebtRemaining = tDomain.TotalAmount
		}
	}

	if amountDebtRemaining.GreaterThan(decimal.Zero) && len(debitedBalanceCategories) == 0 {
		return make(map[int]Category), NewCustomError(CODE_REJECTED_INSUFICIENT_FUNDS, "balance category has insuficient funds")
	}

	return debitedBalanceCategories, nil
}

func (b *Balance) debugLog(msg string) {
	if b.Logger != nil {
		b.Logger.Debug(msg)
	}
}
