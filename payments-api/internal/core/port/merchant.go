package port

type MerchantEntity struct {
	Name          string
	MCC           string
	MappedMccCode string
}

type MerchantRepository interface {
	FindByName(string) (*MerchantEntity, error)
}
