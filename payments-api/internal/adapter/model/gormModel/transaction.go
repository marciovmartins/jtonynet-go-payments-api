package gormModel

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	BaseModel `swaggerignore:"true"`

	UID          uuid.UUID       `json:"uid" binding:"required" example:"91ee2159-f59f-4c89-a543-81987d563d7a" gorm:"type:uuid;unique"`
	AccountID    uint            `json:"account_id" binding:"required" example:"1" gorm:"index:idx_transaction_composite"`
	CategoryID   uint            `json:"category_id" binding:"required" example:"1" gorm:"index:idx_transaction_composite"`
	Amount       decimal.Decimal `json:"amount" binding:"required" example:"110.22" gorm:"type:numeric(20,2);"`
	MCC          string          `json:"mcc" binding:"required" example:"5411" gorm:"type:varchar(5);column:mcc"`
	MerchantName string          `json:"merchant_name" binding:"required" example:"Jonh Doe" gorm:"type:varchar(255)"`

	Category Category `gorm:"foreignKey:CategoryID"`
	Account  Account  `gorm:"foreignKey:AccountID"`
}
