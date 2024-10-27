package gormModel

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Balance struct {
	BaseModel `swaggerignore:"true"`

	UID          uuid.UUID       `json:"uid" binding:"required" example:"4c39c8e9-6278-475a-9513-f2b134ece0b9" gorm:"type:uuid;uniqueIndex"`
	AccountID    uint            `json:"account_id" binding:"required" example:"1"`
	CategoryName string          `json:"category_name" binding:"required" example:"CASH" gorm:"type:varchar(255)"`
	CategoryID   uint            `json:"category_id" binding:"required" example:"1"`
	Amount       decimal.Decimal `json:"amount" binding:"required" example:"110.22" gorm:"type:numeric(20,2);"`

	Account Account `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
