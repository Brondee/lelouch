package storage

import "github.com/Brondee/lelouch/internal/domain"

type Repository interface {
	SaveIfNew(listing domain.Listing) (bool, error)
}
