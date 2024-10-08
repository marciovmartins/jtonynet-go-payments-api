package model

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type GormBaseModel struct {
	gorm.Model
	UUID uuid.UUID `json:"uuid" example:"db047cc5-193a-4989-93f7-08b81c83eea0" type:uuid;unique`
}
