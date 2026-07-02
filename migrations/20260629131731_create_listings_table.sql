-- +goose Up
CREATE TABLE listings (
    id BIGSERIAL PRIMARY KEY,
    external_id TEXT,
    title TEXT NOT NULL,
    description TEXT,
    price BIGINT NOT NULL,
    currency TEXT NOT NULL DEFAULT 'RUB',
    url TEXT NOT NULL,
    image_url TEXT,
    platform TEXT NOT NULL,
    seller_name TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX listings_platform_external_id_idx
ON listings (platform, external_id);


-- +goose Down
DROP INDEX IF EXISTS listings_platform_external_id_idx;
DROP TABLE IF EXISTS listings;
