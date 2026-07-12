package postgres

import (
	"context"
	"fmt"

	"github.com/Brondee/lelouch/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ListingRepository struct {
	DB *pgxpool.Pool
}

func (p *ListingRepository) SaveIfNew(ctx context.Context, listing domain.Listing) (bool, error) {
	const query = `
		INSERT INTO listings (
			external_id,
			title,
			price,
			currency,
			url,
			platform,
			seller_name
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (platform, external_id) DO NOTHING
	`

	tag, err := p.DB.Exec(
		ctx,
		query,
		listing.ID,
		listing.Title,
		listing.Price,
		listing.Currency,
		listing.URL,
		listing.Platform,
		listing.SellerName,
	)
	if err != nil {
		return false, fmt.Errorf("insert listing: %w", err)
	}

	return tag.RowsAffected() == 1, nil
}
