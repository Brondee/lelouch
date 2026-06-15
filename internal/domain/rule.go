package domain

type WatchRule struct {
	ID         string
	Name       string
	MaxPrice   int
	Currency   Currency
	Brands     []string
	Sizes      []string
	SellerName string
	Platform   Platform
	Keywords   []string
	Enabled    bool
}
