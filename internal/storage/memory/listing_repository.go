package memory

import (
	"fmt"

	"github.com/Brondee/lelouch/internal/domain"
)

type ListingRepository struct {
	Listings []domain.Listing
}

func (l *ListingRepository) SaveIfNew(listing domain.Listing) (bool, error) {
	if listing.Title == "" {
		return false, fmt.Errorf("listing is empty")
	}

	for _, repListing := range l.Listings {
		if listing.Title == repListing.Title {
			return false, nil
		}
	}

	l.Listings = append(l.Listings, listing)

	return true, nil
}
