# Lelouch

Lelouch is an early MVP for a future clothing listings parser connected with Telegram bot notifications. Creating project using TDD approach.

## What's done for now:

- Domain models: `Listing`, `WatchRule`, currencies, and platforms
- Filter functions for matching listings by price, platform, brand, size, and keywords. Main function is Matches, which includes all the filter functions
- Money utilities for currency conversion and price normalization
- Unit tests for domain validation, money conversion, and filtering behavior

## Development

Run tests:

```sh
go test ./...