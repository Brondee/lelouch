package parser

import (
	"context"

	"github.com/Brondee/lelouch/internal/domain"
)

type Parser interface {
	Search(ctx context.Context, rule domain.WatchRule) ([]domain.Listing, error)
}
