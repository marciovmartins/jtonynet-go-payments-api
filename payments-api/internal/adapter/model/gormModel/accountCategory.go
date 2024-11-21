package gormModel

type AccountCategory struct {
	BaseModel `swaggerignore:"true"`

	AccountID  uint `json:"account_id" binding:"required" example:"1"`
	CategoryID uint `json:"category_id" binding:"required" example:"1"`

	Account Account `gorm:"foreignKey:AccountID"`
}
