package domain

type Balance struct{}

func (b *Balance) Approve(tDomain Transaction) (Balance, error) {
	return Balance{}, nil
}
