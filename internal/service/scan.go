package service

import (
	"fmt"

	"github.com/Brondee/lelouch/internal/domain"
	"github.com/Brondee/lelouch/internal/filter"
	"github.com/Brondee/lelouch/internal/parser"
	"github.com/Brondee/lelouch/internal/storage"
)

type ScanService struct {
	Repository storage.Repository
	Parser     parser.Parser
}

func (s *ScanService) Scan(rule domain.WatchRule) ([]domain.Listing, error) {
	parsedListings, err := s.Parser.Search()
	if err != nil {
		return nil, fmt.Errorf("parser search: %w", err)
	}

	var matched []domain.Listing

	for _, listing := range parsedListings {
		ok, err := filter.Matches(listing, rule)
		if err != nil {
			return nil, fmt.Errorf("filter matches: %w", err)
		}

		if ok {
			isNew, err := s.Repository.SaveIfNew(listing)
			if err != nil {
				return nil, fmt.Errorf("repository save if new: %w", err)
			}

			if isNew {
				matched = append(matched, listing)
			}
		}
	}

	return matched, nil
}
