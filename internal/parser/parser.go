package parser

import "github.com/Brondee/lelouch/internal/domain"

type Parser interface {
	Search() ([]domain.Listing, error)
}
