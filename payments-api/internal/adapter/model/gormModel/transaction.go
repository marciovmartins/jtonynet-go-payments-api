package gormModel

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	BaseModel `swaggerignore:"true"`

	UID         uuid.UUID       `json:"uid" binding:"required" example:"91ee2159-f59f-4c89-a543-81987d563d7a" gorm:"type:uuid;unique"`
	AccountID   uint            `json:"account_id" binding:"required" example:"1"`
	TotalAmount decimal.Decimal `json:"amount" binding:"required" example:"110.22" gorm:"type:numeric(20,2);"`
	MccCode     string          `json:"mcc_code" binding:"required" example:"5411" gorm:"type:varchar(4);column:mcc_code"`
	Merchant    string          `json:"merchant" binding:"required" example:"Jonh Doe" gorm:"type:varchar(255)"`
}
