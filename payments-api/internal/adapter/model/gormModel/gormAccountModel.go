package model

type GormAccountModel struct {
	GormBaseModel `swaggerignore:"true"`
	Name          string `json:"name" binding:"required" example:"Jonh Doe"`
}
