package gormModel

import "github.com/google/uuid"

type MccCode struct {
	BaseModel `swaggerignore:"true"`

	UID        uuid.UUID `json:"uid" binding:"required" example:"3f77143d-28bb-4d7f-bcf7-0ecff815aab4" gorm:"type:uuid;uniqueIndex"`
	CategoryID uint      `json:"category_id" binding:"required" example:"1"`
	MccCode    string    `json:"mcc_code" binding:"required" example:"5411" gorm:"type:varchar(4);column:mcc_code"`

	Category Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
