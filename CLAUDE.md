# CLAUDE.md

## Project

Go URL shortener service using Gin, PostgreSQL (pgxpool), and a 3-layer architecture (handler → service → repository).

## Commands

- `make run` — Run the server locally
- `make build` — Build the binary
- `make test` — Run unit tests
- `make lint` — Run golangci-lint
- `make docker-up` — Start app + Postgres via Docker Compose
- `make docker-down` — Stop containers
- `DATABASE_URL="postgres://urlshortener:urlshortener@localhost:5432/urlshortener?sslmode=disable" make test` — Run tests including integration tests

## Code Conventions

- 3-layer architecture: handler → service → repository. Keep layers separate.
- Repository layer uses an interface (`URLRepository`) for testability.
- Use `gomock` for mocking interfaces in tests.
- Handler tests use `httptest` with a real Gin router.
- Repository tests are integration tests that require `DATABASE_URL` env var and skip gracefully without it.
- Config is loaded from environment variables, no config files.
- Migrations are plain SQL files in `migrations/` and run automatically on startup.
- Short codes are 7-character random base62 strings.
