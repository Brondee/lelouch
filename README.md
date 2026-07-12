# Lelouch

Lelouch is an early MVP for a future clothing listings parser connected with Telegram bot notifications. Creating project using TDD approach.

## What's done for now:

- Domain models: `Listing`, `WatchRule`, currencies, and platforms
- Filter functions for matching listings by price, platform, brand, size, and keywords. Main function is Matches, which includes all the filter functions
- Money utilities for currency conversion and price normalization
- Unit tests for domain validation, money conversion, and filtering behavior
- Parser interface and FakeParser that imitates real world parser behaviour
- PostgreSQL database in Docker Compose
- Goose migration for the `listings` table
- Repository layer with in-memory and PostgreSQL implementations
- PostgreSQL `ListingRepository` based on `pgx`
- Duplicate detection through `UNIQUE(platform, external_id)` and `ON CONFLICT DO NOTHING`
- Service package that combines all the packages into one Scan pipeline, that takes parsed listings from FakeParser, then filters them via filter package by provided WatchRule and saves new listings through repository storage
- Unit tests for filtering and core `ScanService` scenarios

## Development

Start PostgreSQL:

```sh
docker compose up -d postgres
```

Run migrations:

```sh
goose -dir migrations postgres "postgres://lelouch:12345@localhost:5434/lelouch_db?sslmode=disable" up
```

Run tests:

```sh
go test ./...
```

Run app:

```sh
DATABASE_URL="postgres://lelouch:12345@localhost:5434/lelouch_db?sslmode=disable" go run ./cmd/app
```
