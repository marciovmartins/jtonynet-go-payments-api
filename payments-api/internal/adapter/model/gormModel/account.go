package gormModel

import (
	"github.com/google/uuid"
)

type Account struct {
	BaseModel `swaggerignore:"true"`
	Name      string    `json:"name" binding:"required" example:"Jonh Doe"`
	UID       uuid.UUID `json:"uid" example:"db047cc5-193a-4989-93f7-08b81c83eea0" type:uuid;unique`
}
