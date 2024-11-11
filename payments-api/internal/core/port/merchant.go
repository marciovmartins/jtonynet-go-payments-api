package port

import "context"

type MerchantEntity struct {
	Name string
	MCC  string
}

type MerchantRepository interface {
	FindByName(ctx context.Context, name string) (*MerchantEntity, error)
}
