package fake

import (
	"github.com/Brondee/lelouch/internal/domain"
)

type FakeParser struct {
	Listings []domain.Listing
	Err      error
}

func (f *FakeParser) Search() ([]domain.Listing, error) {
	if f.Err != nil {
		return nil, f.Err
	}

	return f.Listings, nil
}
