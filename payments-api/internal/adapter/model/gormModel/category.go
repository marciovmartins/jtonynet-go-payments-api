package gormModel

import (
	"github.com/google/uuid"
)

type Category struct {
	BaseModel `swaggerignore:"true"`

	UID      uuid.UUID `json:"uid" binding:"required" example:"18e408ce-560a-48a1-b70c-2aa6408b8443" gorm:"type:uuid;uniqueIndex"`
	Name     string    `json:"name" binding:"required" example:"CASH" gorm:"type:varchar(255)"`
	Priority int       `json:"priority" binding:"required" example:"1"`

	MccCodes []MccCode `gorm:"foreignKey:CategoryID"`
	Balances []Balance `gorm:"foreignKey:CategoryID"`
}
