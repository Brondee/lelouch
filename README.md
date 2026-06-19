# Lelouch

Lelouch is an early MVP for a future clothing listings parser connected with Telegram bot notifications. Creating project using TDD approach.

## What's done for now:

- Domain models: `Listing`, `WatchRule`, currencies, and platforms
- Filter functions for matching listings by price, platform, brand, size, and keywords. Main function is Matches, which includes all the filter functions
- Money utilities for currency conversion and price normalization
- Unit tests for domain validation, money conversion, and filtering behavior
- Parser interface and FakeParser that imitates real world parser behaviour
- In-memory ListingRepository with duplicate detection
- Service package that combines all the packages into one Scan pipeline, that takes parsed listings from FakeParser, then filters them via filter package by provided WatchRule and saves them in ListingRepository storage without duplicates
- Unit tests for filtering and core `ScanService` scenarios

## Development

Run tests:

```sh
go test ./...
```

Run app:
```sh
go run ./cmd/app
```