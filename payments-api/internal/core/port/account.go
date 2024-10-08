package port

import model "github.com/jtonynet/go-payments-api/internal/adapter/model/gormModel"

type AccountRepository interface {
	/*
		TODO: USAR DTO!
	*/
	// FindAll() ([]Account, error) //[]domain.Account ???
	FindAll() ([]model.GormAccountModel, error)
}
