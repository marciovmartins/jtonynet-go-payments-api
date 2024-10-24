package port

type MerchantMapEntity struct {
	MerchantName  string
	MccCode       string
	MappedMccCode string
}

type MerchantMapRepository interface {
	FindByMerchantName(string) (MerchantMapEntity, error)
}
