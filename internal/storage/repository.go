package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/Brondee/lelouch/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrorListingExists = errors.New("listing already exists")

type Repository interface {
	SaveIfNew(listing domain.Listing) (bool, error)
}

type PostgresRepository struct {
	Context context.Context
	DB *pgxpool.Pool
}

func (p *PostgresRepository) SaveIfNew(listing domain.Listing) (bool, error) {
	listingDB, err := GetListingByID(p.Context, p.DB, listing.ID)

	if err != nil {
		return false, fmt.Errorf("repository save if new: %w", err)
	}

	if listingDB != nil {
		return false, ErrorListingExists
	}

	return true, nil
}

func GetListingByID(ctx context.Context, pool *pgxpool.Pool, id string) (*domain.Listing, error) {
	const query = `
		SELECT id, title
		FROM listings
		WHERE id = $1
	`

	var listing domain.Listing

	err := pool.QueryRow(ctx, query, id).Scan(
		&listing.ID,
		&listing.Title,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &listing, nil
}
