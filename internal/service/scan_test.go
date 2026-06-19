package service

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Brondee/lelouch/internal/domain"
	"github.com/Brondee/lelouch/internal/money"
	"github.com/Brondee/lelouch/internal/parser"
	"github.com/Brondee/lelouch/internal/parser/fake"
	"github.com/Brondee/lelouch/internal/storage"
	"github.com/Brondee/lelouch/internal/storage/memory"
)

func TestScanService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		rule    domain.WatchRule
		parser  parser.Parser
		storage storage.Repository
		want    []domain.Listing
		wantErr error
	}{
		{
			name:    "returns all listings with empty rule",
			rule:    domain.WatchRule{},
			parser:  &fake.FakeParser{Listings: fakeListings},
			storage: &memory.ListingRepository{},
			want:    fakeListings,
		},
		{
			name:    "returns filtered listings by rule max price",
			rule:    domain.WatchRule{MaxPrice: 35, Currency: domain.USD},
			parser:  &fake.FakeParser{Listings: fakeListings},
			storage: &memory.ListingRepository{},
			want:    append([]domain.Listing{}, fakeListings[0], fakeListings[2]),
		},
		{
			name:    "returns filtered listings by rule keywords",
			rule:    domain.WatchRule{Keywords: []string{"hoodie"}},
			parser:  &fake.FakeParser{Listings: fakeListings},
			storage: &memory.ListingRepository{},
			want:    append([]domain.Listing{}, fakeListings[0]),
		},
		{
			name:    "doesnt save duplicates",
			rule:    domain.WatchRule{},
			parser:  &fake.FakeParser{Listings: fakeListings},
			storage: &memory.ListingRepository{Listings: append([]domain.Listing{}, fakeListings[0])},
			want:    append([]domain.Listing{}, fakeListings[1], fakeListings[2]),
		},
		{
			name:    "returns error when max price is provided and currency is empty",
			rule:    domain.WatchRule{MaxPrice: 35},
			parser:  &fake.FakeParser{Listings: fakeListings},
			storage: &memory.ListingRepository{},
			want:    nil,
			wantErr: money.ErrUnknownCurrency,
		},
		{
			name:    "handles error from parser search",
			rule:    domain.WatchRule{},
			parser:  &fake.FakeParser{Listings: fakeListings, Err: ParserError},
			storage: &memory.ListingRepository{},
			want:    nil,
			wantErr: ParserError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanService := ScanService{Repository: tt.storage, Parser: tt.parser}

			got, err := scanService.Scan(tt.rule)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got error %v, want %v", err, tt.wantErr)
			}
		})
	}
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

var ParserError = errors.New("parser error")
