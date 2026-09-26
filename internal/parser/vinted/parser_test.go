package vinted

import (
	"testing"

	"github.com/Brondee/lelouch/internal/domain"
)

func TestParsePriceCurrency(t *testing.T) {
	tests := []struct {
		name         string
		price        string
		wantPrice    int
		wantCurrency domain.Currency
	}{
		{
			name:         "parses price with decimal comma and euro sign",
			price:        "520,45\u00a0€",
			wantPrice:    521,
			wantCurrency: domain.EUR,
		},
		{
			name:         "parses whole price with euro code",
			price:        "35 EUR",
			wantPrice:    35,
			wantCurrency: domain.EUR,
		},
		{
			name:         "parses rub code",
			price:        "520,00 RUB",
			wantPrice:    520,
			wantCurrency: domain.RUB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrice, gotCurrency, err := parsePriceCurrency(tt.price)
			if err != nil {
				t.Fatalf("parsePriceCurrency() error = %v", err)
			}

			if gotPrice != tt.wantPrice {
				t.Fatalf("price = %d, want %d", gotPrice, tt.wantPrice)
			}

			if gotCurrency != tt.wantCurrency {
				t.Fatalf("currency = %q, want %q", gotCurrency, tt.wantCurrency)
			}
		})
	}
}

func TestParseHTMLToListings(t *testing.T) {
	html := `
		<div data-testid="grid-item">
			<span data-testid="total-combined-price">520,45&nbsp;€</span>
			<p data-testid="catalog-item--description-title">Wool coat</p>
			<p data-testid="catalog-item--description-subtitle">Yohji Yamamoto</p>
		</div>
	`

	listings, err := ParseHTMLToListings([]byte(html))
	if err != nil {
		t.Fatalf("ParseHTMLToListings() error = %v", err)
	}

	if len(listings) != 1 {
		t.Fatalf("len(listings) = %d, want 1", len(listings))
	}

	listing := listings[0]
	if listing.Price != 521 {
		t.Fatalf("listing.Price = %d, want 521", listing.Price)
	}

	if listing.Currency != domain.EUR {
		t.Fatalf("listing.Currency = %q, want %q", listing.Currency, domain.EUR)
	}

	if listing.Title != "Wool coat" {
		t.Fatalf("listing.Title = %q, want %q", listing.Title, "Wool coat")
	}

	if listing.Brand != "Yohji Yamamoto" {
		t.Fatalf("listing.Brand = %q, want %q", listing.Brand, "Yohji Yamamoto")
	}
}
