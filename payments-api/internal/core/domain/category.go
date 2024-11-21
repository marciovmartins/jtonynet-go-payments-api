package domain

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type TransactionCategory struct {
	ID         uint
	CategoryID uint
	Name       string
	Amount     decimal.Decimal
	MCCs       []string
	Priority   int
}

type TransactionByCategories struct {
	Itens map[int]TransactionCategory
}

func (tc *TransactionByCategories) GetByMCC(mcc string) (TransactionCategory, error) {
	for _, transactionCategory := range tc.Itens {
		if searchInSlice(mcc, transactionCategory.MCCs) {
			return transactionCategory, nil
		}
	}

	return TransactionCategory{}, fmt.Errorf("balance category with MCC %s not found", mcc)
}

func (tc *TransactionByCategories) GetFallback() (TransactionCategory, error) {
	var categoryFallback TransactionCategory
	found := false
	maxKey := -1

	for key, transactionCategory := range tc.Itens {
		if key > maxKey && len(transactionCategory.MCCs) == 0 {
			maxKey = key
			categoryFallback = transactionCategory
			found = true
		}
	}

	if found {
		return categoryFallback, nil
	}

	return TransactionCategory{}, fmt.Errorf("balance category with Fallback not found")
}

func (tc *TransactionCategory) Approve(tDomain *Transaction) bool {
	return tc.Amount.Sub(tDomain.Amount).IsPositive()
}

func searchInSlice(search string, slice []string) bool {
	for _, inSlice := range slice {
		if search == inSlice {
			return true
		}
	}
	return false
}
