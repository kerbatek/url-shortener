# URL Shortener

A URL shortener service built with Go, Gin, and PostgreSQL.

## Architecture

```
Handler (Gin) → Service (business logic) → Repository (PostgreSQL)
```

- **Handler**: HTTP request/response handling
- **Service**: URL validation, short code generation (base62)
- **Repository**: CRUD operations via pgxpool

## API

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/shorten` | Create a short URL |
| `GET` | `/:code` | Redirect to original URL |
| `DELETE` | `/url/:id` | Delete a short URL |

### Shorten a URL

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'
```

### Delete a URL

```bash
curl -X DELETE http://localhost:8080/url/1
```

## Running

### Docker (recommended)

```bash
make docker-up
```

Visit `http://localhost:8080` for the web UI.

### Local

Requires a running PostgreSQL instance.

```bash
export APP_PORT=8080
export DB_NAME=urlshortener
export DB_USER=urlshortener
export DB_PASSWORD=urlshortener
export DB_HOST=localhost
export DB_PORT=5432

make run
```

## Development

```bash
make build          # Build binary
make test           # Run unit tests
make lint           # Run linter
make docker-down    # Stop containers
```

### Integration tests

Requires a running PostgreSQL instance:

```bash
DATABASE_URL="postgres://urlshortener:urlshortener@localhost:5432/urlshortener?sslmode=disable" make test
```

## Project Structure

```
cmd/server/          # Application entrypoint
internal/
  handler/           # HTTP handlers (Gin)
  service/           # Business logic
  repository/        # Data access layer
  model/             # Domain models
migrations/          # SQL migration files
static/              # Landing page (HTML/CSS/JS)
```
