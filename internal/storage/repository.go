package storage

import (
	"context"

	"github.com/Brondee/lelouch/internal/domain"
)

type Repository interface {
	SaveIfNew(ctx context.Context, listing domain.Listing) (bool, error)
}
