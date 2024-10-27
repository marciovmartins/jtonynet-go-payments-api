package domain

import (
	"log"

	"github.com/shopspring/decimal"
)

type Balance struct {
	AccountID   uint
	AmountTotal decimal.Decimal
	Categories
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

	bCategoryMCC, err := b.Categories.GetByMCC(tDomain.MccCode)
	if err != nil {
		log.Printf("Error retrieving category by MCC, attempting to debit remaining amount (%s) directly from fallback.\n", amountDebtRemaining)

	} else if bCategoryMCC.Amount.GreaterThanOrEqual(tDomain.TotalAmount) {
		log.Printf("Sufficient funds available in category '%s' to cover the transaction amount.\n", bCategoryMCC.Name)

		amountDebtRemaining = decimal.NewFromFloat(0)

		bCategoryMCC.Amount = bCategoryMCC.Amount.Sub(tDomain.TotalAmount)
		debitedBalanceCategories[bCategoryMCC.Priority] = bCategoryMCC

	} else if bCategoryMCC.Amount.IsPositive() {
		log.Printf("Category '%s' has a positive balance but insufficient funds for the transaction. Debiting the entire amount and checking fallback category.\n", bCategoryMCC.Name)

		amountDebtRemaining = amountDebtRemaining.Sub(bCategoryMCC.Amount)

		bCategoryMCC.Amount = decimal.NewFromFloat(0)
		debitedBalanceCategories[bCategoryMCC.Priority] = bCategoryMCC
	}

	if amountDebtRemaining.GreaterThan(decimal.Zero) {
		bCategoryFallback, err := b.Categories.GetFallback()
		if err != nil {
			log.Println("Error retrieving fallback category. Rolling back changes and returning an error.")

			return make(map[int]Category), NewCustomError(CODE_REJECTED_GENERIC, "Category Fallback not found")

		} else if bCategoryFallback.Amount.GreaterThanOrEqual(amountDebtRemaining) {
			log.Println("Fallback category has sufficient funds available for the transaction.")

			bCategoryFallback.Amount = bCategoryFallback.Amount.Sub(amountDebtRemaining)
			debitedBalanceCategories[bCategoryFallback.Priority] = bCategoryFallback

			amountDebtRemaining = decimal.NewFromFloat(0)
		} else {
			log.Println("Insufficient funds in the fallback category. Rolling back changes and returning an error.")

			debitedBalanceCategories = make(map[int]Category)
			amountDebtRemaining = tDomain.TotalAmount
		}
	}

	if amountDebtRemaining.GreaterThan(decimal.Zero) && len(debitedBalanceCategories) == 0 {
		return make(map[int]Category), NewCustomError(CODE_REJECTED_INSUFICIENT_FUNDS, "balance category has insuficient funds")
	}

	return debitedBalanceCategories, nil
}
