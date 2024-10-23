package port

type MerchantMapEntity struct {
	MerchantName  string
	MccCode       string
	MappedMccCode string
}

type MerchantMaptRepository interface {
	FindByMerchantName(string) (MerchantMapEntity, error)
}
