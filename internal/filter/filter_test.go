package filter

import (
	"errors"
	"testing"

	"github.com/Brondee/lelouch/internal/domain"
	"github.com/Brondee/lelouch/internal/money"
)

func TestMatchesPrice(t *testing.T) {
	t.Parallel()

	matchesPriceTests := []struct {
		name    string
		listing domain.Listing
		rule    domain.WatchRule
		want    bool
		wantErr error
	}{
		{
			name:    "matches price in the same currency",
			listing: domain.Listing{Price: 50, Currency: domain.USD},
			rule:    domain.WatchRule{MaxPrice: 100, Currency: domain.USD},
			want:    true,
		},
		{
			name:    "matches price in a different currency",
			listing: domain.Listing{Price: 50, Currency: domain.USD},
			rule:    domain.WatchRule{MaxPrice: 100, Currency: domain.RUB},
			want:    false,
		},
		{
			name:    "returns error when listing currency is unknown",
			listing: domain.Listing{Price: 50, Currency: "HUI"},
			rule:    domain.WatchRule{MaxPrice: 100, Currency: domain.USD},
			want:    false,
			wantErr: money.ErrUnknownCurrency,
		},
		{
			name:    "returns error when rule currency is unknown",
			listing: domain.Listing{Price: 50, Currency: domain.USD},
			rule:    domain.WatchRule{MaxPrice: 100, Currency: "HUI"},
			want:    false,
			wantErr: money.ErrUnknownCurrency,
		},
		{
			name:    "returns error when price is negative",
			listing: domain.Listing{Price: -50, Currency: domain.USD},
			rule:    domain.WatchRule{MaxPrice: 100, Currency: domain.EUR},
			want:    false,
			wantErr: money.ErrNegativePrice,
		},
	}

	for _, tt := range matchesPriceTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MatchesPrice(tt.listing, tt.rule)

			if got != tt.want {
				t.Errorf("got %v want %v", got, tt.want)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got error %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestMatchesPlatform(t *testing.T) {
	t.Parallel()

	matchesPlatformTests := []struct {
		name    string
		listing domain.Listing
		rule    domain.WatchRule
		want    bool
		wantErr error
	}{
		{
			name:    "matches same platforms",
			listing: domain.Listing{Platform: domain.PlatformVinted},
			rule:    domain.WatchRule{Platform: domain.PlatformVinted},
			want:    true,
		},
		{
			name:    "doesnt match different platforms",
			listing: domain.Listing{Platform: domain.PlatformVinted},
			rule:    domain.WatchRule{Platform: domain.PlatformMercari},
			want:    false,
		},
		{
			name:    "returns error when trying to match unknown listing platform",
			listing: domain.Listing{Platform: "HUI"},
			rule:    domain.WatchRule{Platform: domain.PlatformVinted},
			want:    false,
			wantErr: domain.ErrInvalidPlatform,
		},
		{
			name:    "returns error when trying to match unknown rule platform",
			listing: domain.Listing{Platform: domain.PlatformVinted},
			rule:    domain.WatchRule{Platform: "HUI"},
			want:    false,
			wantErr: domain.ErrInvalidPlatform,
		},
	}

	for _, tt := range matchesPlatformTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MatchesPlatform(tt.listing, tt.rule)

			if got != tt.want {
				t.Errorf("got %v want %v", got, tt.want)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got error %v want %v", err, tt.wantErr)
			}
		})
	}
}

func TestMatchesSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		listing domain.Listing
		rule    domain.WatchRule
		want    bool
	}{
		{
			name:    "matches exact size",
			listing: domain.Listing{Size: "m"},
			rule:    domain.WatchRule{Sizes: []string{"m", "l"}},
			want:    true,
		},
		{
			name:    "matches size ignoring case",
			listing: domain.Listing{Size: "m"},
			rule:    domain.WatchRule{Sizes: []string{"M", "l"}},
			want:    true,
		},
		{
			name:    "matches size with spaces",
			listing: domain.Listing{Size: " M "},
			rule:    domain.WatchRule{Sizes: []string{"M", "l"}},
			want:    true,
		},
		{
			name:    "matches size when rule slice of sizes is empty",
			listing: domain.Listing{Size: "xl"},
			rule:    domain.WatchRule{},
			want:    true,
		},
		{
			name:    "doesnt match different sizes",
			listing: domain.Listing{Size: "xl"},
			rule:    domain.WatchRule{Sizes: []string{"m", "l"}},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchesSize(tt.listing, tt.rule)

			if got != tt.want {
				t.Errorf("got %v want %v", got, tt.want)
			}
		})
	}
}

func TestMatchesBrand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		listing domain.Listing
		rule    domain.WatchRule
		want    bool
	}{
		{
			name:    "matches exact brand",
			listing: domain.Listing{Brand: "y3"},
			rule:    domain.WatchRule{Brands: []string{"y3", "yohji yamamoto"}},
			want:    true,
		},
		{
			name:    "matches brand ignoring case",
			listing: domain.Listing{Brand: "Y3"},
			rule:    domain.WatchRule{Brands: []string{"y3", "yohji yamamoto"}},
			want:    true,
		},
		{
			name:    "matches brand ignoring spaces",
			listing: domain.Listing{Brand: " y3 "},
			rule:    domain.WatchRule{Brands: []string{"y3", "yohji yamamoto"}},
			want:    true,
		},
		{
			name:    "matches brand when rule slices of brands is empty",
			listing: domain.Listing{Brand: "y3"},
			rule:    domain.WatchRule{},
			want:    true,
		},
		{
			name:    "doesnt match different brands",
			listing: domain.Listing{Brand: "maison margiela"},
			rule:    domain.WatchRule{Brands: []string{"y3", "yohji yamamoto"}},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchesBrand(tt.listing, tt.rule)

			if got != tt.want {
				t.Errorf("got %v want %v", got, tt.want)
			}
		})
	}
}

func TestMatchesKeywords(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		listing domain.Listing
		rule    domain.WatchRule
		want    bool
	}{
		{
			name:    "matches exact keyword in listing title",
			listing: domain.Listing{Title: "black hoodie archive"},
			rule:    domain.WatchRule{Keywords: []string{"archive"}},
			want:    true,
		},
		{
			name:    "matches keyword in listing title ignoring case",
			listing: domain.Listing{Title: "black hoodie archive"},
			rule:    domain.WatchRule{Keywords: []string{"Archive"}},
			want:    true,
		},
		{
			name:    "matches keyword phrase in listing title ignoring case",
			listing: domain.Listing{Title: "black hoodie archive"},
			rule:    domain.WatchRule{Keywords: []string{"hoodie archive"}},
			want:    true,
		},
		{
			name:    "matches keyword in listing when rule slice of keywords is empty",
			listing: domain.Listing{Title: "black hoodie archive"},
			rule:    domain.WatchRule{},
			want:    true,
		},
		{
			name:    "doesnt match different keywords",
			listing: domain.Listing{Title: "black hoodie minimalistic"},
			rule:    domain.WatchRule{Keywords: []string{"archive"}},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchesKeywords(tt.listing, tt.rule)

			if got != tt.want {
				t.Errorf("got %v want %v", got, tt.want)
			}
		})
	}
}

func TestMatches(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		listing   domain.Listing
		rule      domain.WatchRule
		want      bool
		wantError error
	}{
		{
			name:    "matches full listing params",
			listing: validListing(),
			rule:    validRule(),
			want:    true,
		},
		{
			name: "doesnt match price",
			listing: listingWith(func(listing *domain.Listing) {
				listing.Price = 10000
			}),
			rule: validRule(),
			want: false,
		},
		{
			name: "returns error when currency is unknown",
			listing: listingWith(func(listing *domain.Listing) {
				listing.Currency = "HUI"
			}),
			rule:      validRule(),
			want:      false,
			wantError: money.ErrUnknownCurrency,
		},
		{
			name: "doesnt match platform",
			listing: listingWith(func(listing *domain.Listing) {
				listing.Platform = domain.PlatformAvito
			}),
			rule: validRule(),
			want: false,
		},
		{
			name: "returns error when platform is invalid",
			listing: listingWith(func(listing *domain.Listing) {
				listing.Platform = "HUI"
			}),
			rule:      validRule(),
			want:      false,
			wantError: domain.ErrInvalidPlatform,
		},
		{
			name: "doesnt match size",
			listing: listingWith(func(listing *domain.Listing) {
				listing.Size = "xl"
			}),
			rule: validRule(),
			want: false,
		},
		{
			name: "doesnt match brand",
			listing: listingWith(func(listing *domain.Listing) {
				listing.Brand = "marcelo miracles"
			}),
			rule: validRule(),
			want: false,
		},
		{
			name: "doesnt match keywords",
			listing: listingWith(func(listing *domain.Listing) {
				listing.Title = "cool tshirt bbc"
			}),
			rule: validRule(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Matches(tt.listing, tt.rule)

			if got != tt.want {
				t.Errorf("got %v want %v", got, tt.want)
			}

			if !errors.Is(err, tt.wantError) {
				t.Fatalf("got error %v want %v", err, tt.wantError)
			}
		})
	}
}

func listingWith(change func(*domain.Listing)) domain.Listing {
	listing := validListing()
	change(&listing)
	return listing
}

func validListing() domain.Listing {
	return domain.Listing{
		Title:      "black hoodie archive",
		Platform:   domain.PlatformVinted,
		Price:      100,
		Currency:   domain.EUR,
		Size:       "m",
		Brand:      "y3",
		SellerName: "anton",
		Category:   "hoodies",
		URL:        "vinted.de/hoodie",
	}
}

func validRule() domain.WatchRule {
	return domain.WatchRule{
		MaxPrice:   120,
		Currency:   domain.EUR,
		Brands:     []string{"y3", "maison margiela", "hysteric glamour"},
		Sizes:      []string{"m", "l"},
		Platform:   domain.PlatformVinted,
		Keywords:   []string{"hoodie", "archive", "old", "2000s"},
		Enabled:    true,
		Name:       "hoodie rule",
		SellerName: "anton",
	}
}
