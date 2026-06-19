package filter

import (
	"fmt"
	"strings"

	"github.com/Brondee/lelouch/internal/domain"
	"github.com/Brondee/lelouch/internal/money"
)

func MatchesPrice(listing domain.Listing, rule domain.WatchRule) (bool, error) {
	if rule.MaxPrice == 0 {
		return true, nil
	}

	curPriceUsd, err := money.ToUSD(listing.Price, listing.Currency)

	if err != nil {
		return false, fmt.Errorf("listing price: %w", err)
	}

	maxPriceUsd, err := money.ToUSD(rule.MaxPrice, rule.Currency)

	if err != nil {
		return false, fmt.Errorf("rule max price: %w", err)
	}

	return curPriceUsd <= maxPriceUsd, nil
}

func MatchesPlatform(listing domain.Listing, rule domain.WatchRule) (bool, error) {
	if rule.Platform == "" {
		return true, nil
	}

	err := listing.Platform.Validate()
	if err != nil {
		return false, fmt.Errorf("listing platform: %w", err)
	}

	err = rule.Platform.Validate()
	if err != nil {
		return false, fmt.Errorf("rule platform: %w", err)
	}

	return listing.Platform == rule.Platform, nil
}

func MatchesSize(listing domain.Listing, rule domain.WatchRule) bool {
	if len(rule.Sizes) == 0 {
		return true
	}

	listingSize := normalizeText(listing.Size)

	for _, size := range rule.Sizes {
		if listingSize == normalizeText(size) {
			return true
		}
	}

	return false
}

func MatchesBrand(listing domain.Listing, rule domain.WatchRule) bool {
	if len(rule.Brands) == 0 {
		return true
	}

	listingBrand := normalizeText(listing.Brand)

	for _, brand := range rule.Brands {
		if listingBrand == normalizeText(brand) {
			return true
		}
	}

	return false
}

func MatchesKeywords(listing domain.Listing, rule domain.WatchRule) bool {
	if len(rule.Keywords) == 0 {
		return true
	}

	normalizedTitle := normalizeText(listing.Title)

	for _, keyword := range rule.Keywords {
		normalizedKeyword := normalizeText(keyword)
		if normalizedKeyword == "" {
			continue
		}

		if strings.Contains(normalizedTitle, normalizedKeyword) {
			return true
		}
	}

	return false
}

func Matches(listing domain.Listing, rule domain.WatchRule) (bool, error) {
	matchesPrice, priceErr := MatchesPrice(listing, rule)
	if priceErr != nil {
		return false, fmt.Errorf("listing price: %w", priceErr)
	}

	if !matchesPrice {
		return false, nil
	}

	matchesPlatform, platformErr := MatchesPlatform(listing, rule)
	if platformErr != nil {
		return false, fmt.Errorf("listing platform: %w", platformErr)
	}

	if !matchesPlatform {
		return false, nil
	}

	if !MatchesSize(listing, rule) {
		return false, nil
	}

	if !MatchesBrand(listing, rule) {
		return false, nil
	}

	if !MatchesKeywords(listing, rule) {
		return false, nil
	}

	return true, nil
}

func normalizeText(text string) string {
	return strings.ToLower(strings.TrimSpace(text))
}
