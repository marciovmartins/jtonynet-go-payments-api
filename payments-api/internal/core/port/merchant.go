package port

type MerchantEntity struct {
	Name          string
	MccCode       string
	MappedMccCode string
}

type MerchantRepository interface {
	FindByName(string) (*MerchantEntity, error)
}
