package gormModel

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Balance struct {
	BaseModel    `swaggerignore:"true"`
	UID          uuid.UUID       `json:"uid" example:"91ee2159-f59f-4c89-a543-81987d563d7a" type:uuid;unique`
	AccountID    uint            `json:"account_id" binding:"required" example:"1" type:uint`
	CategoryName string          `json:"category_name" binding:"required" example:"CASH"`
	Amount       decimal.Decimal `json:"amount" binding:"required" example:"100.00" sql:"type:numeric(20,2);"`
}
