package port

type MerchantEntity struct {
	Name string
	MCC  string
}

type MerchantRepository interface {
	FindByName(string) (*MerchantEntity, error)
}
