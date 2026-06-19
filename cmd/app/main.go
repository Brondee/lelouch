package main

import (
	"fmt"
	"os"

	"github.com/Brondee/lelouch/internal/domain"
	"github.com/Brondee/lelouch/internal/parser/fake"
	"github.com/Brondee/lelouch/internal/service"
	"github.com/Brondee/lelouch/internal/storage/memory"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	parser := &fake.FakeParser{Listings: fakeListings}
	storage := &memory.ListingRepository{}

	scanService := service.ScanService{Repository: storage, Parser: parser}

	ruleByMaxPrice := domain.WatchRule{MaxPrice: 35, Currency: domain.USD}

	listings, err := scanService.Scan(ruleByMaxPrice)
	if err != nil {
		return fmt.Errorf("scan listings: %w", err)
	}

	for _, listing := range listings {
		fmt.Println(listing)
	}

	return nil
}

var fakeListings = []domain.Listing{
	{
		ID:         "123",
		Title:      "y3 hoodie",
		Platform:   domain.PlatformMercari,
		SellerName: "ivan",
		Brand:      "y3",
		Price:      34,
		Currency:   domain.USD,
		Size:       "m",
		Category:   "hoodie",
		URL:        "https://mercari.jp/y3-hoodie",
	},
	{
		ID:         "1234",
		Title:      "number nine tshirt",
		Platform:   domain.PlatformVinted,
		SellerName: "anton",
		Brand:      "number nine",
		Price:      45,
		Currency:   domain.EUR,
		Size:       "l",
		Category:   "tshirt",
		URL:        "https://vinted.de/number-nine-tshirt",
	},
	{
		ID:         "12345",
		Title:      "oakley long sleeve",
		Platform:   domain.PlatformVinted,
		SellerName: "gregori",
		Brand:      "oakley",
		Price:      23,
		Currency:   domain.EUR,
		Size:       "l",
		Category:   "long sleeve",
		URL:        "https://vinted.de/oakley-long-sleeve",
	},
}
