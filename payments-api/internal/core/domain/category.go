package domain

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Category struct {
	ID       uint
	Name     string
	Amount   decimal.Decimal
	MCCs     []string
	Priority int
}

type Categories struct {
	Itens map[int]Category
}

func (c *Categories) GetByMCC(mcc string) (Category, error) {
	for _, categoryItem := range c.Itens {
		if searchInSlice(mcc, categoryItem.MCCs) {
			return categoryItem, nil
		}
	}

	return Category{}, fmt.Errorf("balance category with MCC %s not found", mcc)
}

func (c *Categories) GetFallback() (Category, error) {
	var categoryFallback Category
	found := false
	maxKey := -1

	for key, categoryItem := range c.Itens {
		if key > maxKey && len(categoryItem.MCCs) == 0 {
			maxKey = key
			categoryFallback = categoryItem
			found = true
		}
	}

	if found {
		return categoryFallback, nil
	}

	return Category{}, fmt.Errorf("balance category with Fallback not found")
}

func (c *Category) Approve(tDomain *Transaction) bool {
	return c.Amount.Sub(tDomain.TotalAmount).IsPositive()
}

func searchInSlice(search string, slice []string) bool {
	for _, inSlice := range slice {
		if search == inSlice {
			return true
		}
	}
	return false
}
