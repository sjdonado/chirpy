# Chirpy API

Chirpy is a small Twitter-like HTTP API built for learning: users post chirps (tweets). It follows the Boot.dev “Learn HTTP Servers in Go” course, with my own touches around code organization, folder structure, routing, and DRY patterns.

## Tech Stack
- Go + chi router for minimal routing capabilities
- PostgreSQL 15 (Docker)
- sqlc for typed DB access
- JWT access tokens + refresh tokens
- argon2id password hashing

## Project Structure
- `main.go` — router wiring and server setup
- `handler/` — HTTP handlers
- `middleware/` — auth + metrics middleware
- `internal/api/` — config, headers, and response helpers
- `internal/auth/` — JWT + password helpers
- `internal/database/` — sqlc-generated queries/models
- `sql/schema/` — SQL schema (migrations)
- `sql/queries/` — SQL queries for sqlc

## Getting Started (Local)

### Prerequisites
- Go (see `go.mod`)
- Docker Desktop (for Postgres)

### 1) Start Postgres with Docker
```
docker compose up -d
```

To stop and remove the container (and data volume):
```
docker compose down -v
```

### 2) Configure environment
Create a `.env` (or export env vars directly):

```
DB_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
PLATFORM=dev
JWT_SECRET=your-base64-or-random-secret
POLKA_KEY=your-polka-webhook-key
```

### 3) Apply schema
Run the SQL files in `sql/schema/` in order (001 → 005) using your preferred migration tool, or use `goose`:

```
goose postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up
```

### 4) Run the server
```
go run ./main.go
```

Server starts on `http://localhost:8080`.

## Routes

### Public
- `GET /api/healthz`
- `POST /api/users` — create user
- `POST /api/login` — returns access + refresh tokens
- `POST /api/refresh` — uses refresh token (Authorization: Bearer)
- `POST /api/revoke` — revoke refresh token (Authorization: Bearer)
- `GET /api/chirps` — supports `author_id` and `sort=asc|desc`
- `GET /api/chirps/{id}`

### Authenticated (JWT)
- `PUT /api/users`
- `POST /api/chirps`
- `DELETE /api/chirps/{id}`

Auth header:
```
Authorization: Bearer <access_token>
```

### Webhook
- `POST /api/polka/webhooks`

Webhook auth header:
```
Authorization: ApiKey <POLKA_KEY>
```

### Admin
- `GET /admin/metrics` — HTML metrics page
- `POST /admin/reset` — resets metrics (and deletes users when `PLATFORM=dev`)

### Static
- `GET /app/*` — serves files from the repo root

## Notes
This project is intentionally small and focused on learning: clean handlers, minimal middleware, and sqlc-driven queries for type safety.
