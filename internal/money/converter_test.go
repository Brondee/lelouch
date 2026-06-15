package money

import (
	"testing"

	"github.com/Brondee/lelouch/internal/domain"
)

func TestToUSD(t *testing.T) {
	toUSDTests := []struct {
		name    string
		listing domain.Listing
		want    float64
		wantErr error
	}{
		{
			name:    "convert to a known currency",
			listing: domain.Listing{Price: 100, Currency: domain.EUR},
			want:    116.0,
		},
		{
			name:    "convert to a unknown currency",
			listing: domain.Listing{Price: 100, Currency: "HUI"},
			want:    0,
			wantErr: ErrUnknownCurrency,
		},
		{
			name:    "returns error when a negative price passed",
			listing: domain.Listing{Price: -100, Currency: domain.EUR},
			want:    0,
			wantErr: ErrNegativePrice,
		},
	}

	for _, tt := range toUSDTests {
		got, err := ToUSD(tt.listing.Price, tt.listing.Currency)

		if got != tt.want {
			t.Errorf("got %v want %v", got, tt.want)
		}

		if err != tt.wantErr {
			t.Fatalf("got error %v, want %v", err, tt.wantErr)
		}
	}
}
