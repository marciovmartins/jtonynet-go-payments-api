package domain

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Category struct {
	ID       int
	Name     string
	Amount   decimal.Decimal
	MCCcodes []string
	Order    int
}

type Categories struct {
	Itens map[int]Category
}

func (c *Categories) GetByMCC(mcc string) (Category, error) {
	for _, categoryItem := range c.Itens {
		if searchInSlice(mcc, categoryItem.MCCcodes) {
			return categoryItem, nil
		}
	}

	return Category{}, fmt.Errorf("balance category with MCC %s not found", mcc)
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
