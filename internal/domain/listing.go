package domain

type Listing struct {
	ID         string
	Title      string
	Platform   Platform
	SellerName string
	Brand      string
	Price      int
	Currency   Currency
	Size       string
	Category   string
	URL        string
}
