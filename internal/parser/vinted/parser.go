package vinted

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Brondee/lelouch/internal/domain"
	"github.com/PuerkitoBio/goquery"
)

var vintedBaseUrl = "https://www.vinted.de"
var orderByNewest = "&order=newest_first"

type Parser struct {
	client *Client
}

func NewParser(client *Client) *Parser {
	return &Parser{client: client}
}

func (p *Parser) Search(ctx context.Context, rule domain.WatchRule) ([]domain.Listing, error) {
	separatedBrandName := strings.Fields(rule.Brands[0])

	normalizedBrandName := ""
	for index, word := range separatedBrandName {
		if index != len(word) {
			normalizedBrandName += word + "%20"
		}
	}

	url := vintedBaseUrl + "/catalog?search_text=" + normalizedBrandName + orderByNewest

	body, err := p.client.Fetch(url)
	if err != nil {
		return nil, fmt.Errorf("fetch url: %w", err)
	}

	listings, err := ParseHTMLToListings(body)
	if err != nil {
		return nil, fmt.Errorf("parse html to listings: %w", err)
	}

	return listings, nil
}

func ParseHTMLToListings(body []byte) ([]domain.Listing, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	var listings []domain.Listing

	var parseErr error

	doc.Find(`[data-testid="grid-item"]`).EachWithBreak(func(i int, item *goquery.Selection) bool {
		price := strings.TrimSpace(
			item.Find(`[data-testid="total-combined-price"]`).First().Text(),
		)
		priceInt, currency, err := parsePriceCurrency(price)
		if err != nil {
			parseErr = fmt.Errorf("parse listing %d price %q: %w", i, price, err)
			return false
		}

		overlay := item.Find(`a[data-testid^="product-item-id-"][data-testid$="--overlay-link"]`).First()

		href, ok := overlay.Attr("href")
		if !ok {
			panic("error getting href attr from overlay")
		}

		titleAttr, ok := overlay.Attr("title")
		if !ok {
			panic("error getting title attr from overlay")
		}

		info, err := parseVintedOverlayTitle(titleAttr)
		if err != nil {
			fmt.Print(titleAttr)
			panic(err)
		}

		listing := domain.Listing{
			Price:    priceInt,
			Currency: currency,
			Title:    info[0],
			Brand:    info[1],
			Size:     info[3],
			URL:      href,
		}

		listings = append(listings, listing)

		return true
	})

	if parseErr != nil {
		return nil, parseErr
	}

	return listings, nil
}

func parsePriceCurrency(rawPrice string) (int, domain.Currency, error) {
	parts := strings.Fields(rawPrice)
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("expected price and currency")
	}

	amount, err := strconv.ParseFloat(strings.ReplaceAll(parts[0], ",", "."), 64)
	if err != nil {
		return 0, "", fmt.Errorf("parse price amount: %w", err)
	}

	currency, err := parseCurrency(parts[1])
	if err != nil {
		return 0, "", err
	}

	return int(math.Ceil(amount)), currency, nil
}

func parseCurrency(rawCurrency string) (domain.Currency, error) {
	switch strings.ToUpper(strings.TrimSpace(rawCurrency)) {
	case "€", "EUR":
		return domain.EUR, nil
	case "RUB":
		return domain.RUB, nil
	case "$", "USD":
		return domain.USD, nil
	default:
		return "", fmt.Errorf("unknown currency %q", rawCurrency)
	}
}

func parseVintedOverlayTitle(data string) ([]string, error) {
	var result []string

	titleWithBrand := strings.Split(data, ", brand:")
	if len(titleWithBrand) < 2 {
		return nil, fmt.Errorf("didnt get title or brand from overlay title")
	}

	title := titleWithBrand[0]
	result = append(result, title)

	brandWithCondition := strings.Split(titleWithBrand[1], ", condition:")
	if len(brandWithCondition) < 2 {
		return nil, fmt.Errorf("didnt get brand or condition from overlay title")
	}

	brand := brandWithCondition[0]
	result = append(result, brand)

	conditionWithSize := strings.Split(brandWithCondition[1], ", size:")
	if len(conditionWithSize) < 2 {
		return nil, fmt.Errorf("didnt get condition or size from overlay title")
	}

	condition := conditionWithSize[0]
	result = append(result, condition)

	size := strings.Split(conditionWithSize[1], ", ")
	if size[0] == "" {
		return nil, fmt.Errorf("didnt get size from overlay title")
	}

	result = append(result, size[0])

	return result, nil
}
