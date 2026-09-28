# Backend Service (Go + PostgreSQL + pgx + SQLC)

## 1. Overview & Technology Stack

This service is an idiomatic, production-grade Go HTTP REST API server with PostgreSQL persistence.

* **Go Version**: `1.23+` (configured in `go.mod`, verified on stable `go1.27.1`).
* **Routing**: Standard library `net/http` with Go 1.22+ enhanced ServeMux method matching.
* **Database Driver**: `github.com/jackc/pgx/v5` and `pgxpool`.
* **Database Query Engine**: `sqlc` for compile-safe Go query code generation.
* **Database Migrations**: Embedded transactional SQL migrations.
* **Logging**: Standard library `log/slog` for structured logging.
* **Testing**: Standard `testing` and `httptest` packages.

---

## 2. Package Architecture

```text
backend/
├── cmd/
│   └── server/
│       └── main.go              # Minimal bootstrap & graceful shutdown entrypoint
├── internal/                    # Private implementation packages
│   ├── config/                  # Environment loading & validation (only place os.Getenv is allowed)
│   ├── database/                # Database connection pool and SQL persistence
│   │   ├── migrations/          # Version-controlled SQL schema changes
│   │   ├── queries/             # Hand-authored SQL queries for SQLC
│   │   ├── sqlc/                # Auto-generated type-safe Go SQL code (DO NOT EDIT)
│   │   ├── pool.go              # pgxpool lifecycle and health check
│   │   └── migrate.go           # Transactional embedded SQL migration runner
│   ├── health/                  # Health check domain module
│   │   ├── handler.go           # HTTP endpoints for /health, /live, /ready
│   │   ├── service.go           # Health and readiness evaluation
│   │   ├── types.go             # Domain health response structs
│   │   └── handler_test.go      # httptest unit tests
│   ├── http/                    # HTTP transport layer
│   │   ├── middleware/          # Request ID, slog logger, recoverer, CORS, security, limits
│   │   ├── response/            # Standardized JSON response envelopes & error formatting
│   │   └── routes/              # Route registration and middleware pipeline assembly
│   └── logger/                  # Structured slog initialization and context propagation
├── sqlc.yaml                    # SQLC configuration for pgx/v5
├── go.mod                       # Go module dependencies
├── go.sum                       # Go module checksums
├── .env.example                 # Environment variables template
├── Dockerfile                   # Multi-stage production container
└── README.md                    # This document
```

---

## 3. Local Development

### 3.1 Prerequisites
* Go 1.23+ installed.
* PostgreSQL running locally or via Docker:
  ```bash
  # From project root:
  docker compose up -d postgres
  ```

### 3.2 Running the Backend
```bash
cd backend
cp .env.example .env
go mod download
go run cmd/server/main.go
```
The server will bind to `http://localhost:8080`.

---

## 4. Environment Variables

| Variable | Default Value | Description |
| :--- | :--- | :--- |
| `APP_ENV` | `development` | Runtime environment (`development`, `staging`, `production`) |
| `HTTP_PORT` | `8080` | Port for the HTTP server to bind to |
| `DATABASE_URL` | `postgres://...` | Connection URI for PostgreSQL |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:5173,http://localhost:3000` | Comma-separated allowed CORS origins |
| `READ_TIMEOUT_SECONDS` | `10` | Maximum duration for reading request |
| `WRITE_TIMEOUT_SECONDS` | `10` | Maximum duration for writing response |
| `IDLE_TIMEOUT_SECONDS` | `120` | Maximum duration to keep idle connections open |
| `SHUTDOWN_TIMEOUT_SECONDS`| `15` | Maximum grace period during shutdown |
| `DB_MAX_OPEN_CONNS` | `25` | Maximum open connections in pgxpool |
| `DB_MIN_OPEN_CONNS` | `5` | Minimum idle connections in pgxpool |
| `DB_MAX_IDLE_TIME_MINUTES`| `15` | Idle connection timeout |
| `LOG_LEVEL` | `info` | Minimum log level (`debug`, `info`, `warn`, `error`) |
| `LOG_FORMAT` | `text` | Log format (`text` or `json`) |

---

## 5. SQLC Generation & Migrations

### Running Migrations
Migrations run automatically upon application startup using Go's `embed` package. To add a new migration:
1. Create `internal/database/migrations/<sequence>_<name>.up.sql`.
2. Create `internal/database/migrations/<sequence>_<name>.down.sql`.

### Generating SQLC Code
```bash
sqlc generate
```
Reads `internal/database/queries/*.sql` and schema files to regenerate `internal/database/sqlc/*.go`.

---

## 6. Testing, Race Detection & Linting

```bash
# Run unit and handler tests
go test -v ./...

# Run tests with the Go race detector
go test -race ./...

# Format all Go source files
gofmt -s -w .

# Run static analysis
go vet ./...

# Build binary
go build -o bin/server ./cmd/server/main.go
```

---

## 7. Graceful Shutdown & Production Deployment

On receiving `SIGINT` or `SIGTERM`:
1. Server stops accepting new inbound connections.
2. In-flight requests are granted up to `SHUTDOWN_TIMEOUT_SECONDS` (15s default) to complete.
3. Database connection pool is closed cleanly via `pool.Close()`.
4. Process exits with code `0`.
