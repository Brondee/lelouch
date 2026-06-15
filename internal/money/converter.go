package money

import (
	"errors"

	"github.com/Brondee/lelouch/internal/domain"
)

var (
	ErrUnknownCurrency = errors.New("unknown currency")
	ErrNegativePrice   = errors.New("negative price")
)

var currencyToUSD = map[domain.Currency]float64{
	domain.USD: 1,
	domain.EUR: 1.16,
	domain.RUB: 0.014,
}

func ToUSD(price int, currency domain.Currency) (float64, error) {
	if price < 0 {
		return 0, ErrNegativePrice
	}

	rate, ok := currencyToUSD[currency]

	if !ok {
		return 0, ErrUnknownCurrency
	}

	converted := float64(price) * rate
	return roundToTwo(converted), nil
}
