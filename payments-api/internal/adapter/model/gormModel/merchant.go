package gormModel

import (
	"github.com/google/uuid"
)

type Merchant struct {
	BaseModel `swaggerignore:"true"`
	UID       uuid.UUID `json:"uid" binding:"required" example:"0b0364a1-4955-48b6-8c63-8a446b918682" gorm:"type:uuid;unique"`
	Name      string    `json:"name" binding:"required" example:"UBER EATS   SAO PAULO BR" gorm:"type:varchar(255);index"`
	MccCode   string    `json:"mcc_code" binding:"required" example:"5555" gorm:"type:varchar(255);column:mcc_code"`
	MCC       string    `json:"mcc" binding:"required" example:"5411" gorm:"type:varchar(5);column:mcc"`
	MccID     uint      `json:"mcc_id" binding:"required" example:"1"`
}
