# Local Development Guide

## 1. Prerequisites and Tooling

Ensure the following tools are installed on your host system:

* **Go**: Version `1.23` or higher (verified with stable `go1.27.1` / `go1.23+`).
* **Node.js**: Version `20.x` or higher (LTS recommended; tested with `v24.x`).
* **npm**: Version `10.x` or higher.
* **PostgreSQL**: Version `15` or `16` (or run via Docker Compose).
* **Docker & Docker Compose**: For containerized development and CI builds.
* **SQLC** (Optional for local query regeneration): Installed via `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`.

---

## 2. Quickstart with Docker Compose

The fastest way to spin up the entire isolated stack (PostgreSQL, Go API server, and React frontend) is with Docker Compose:

```bash
# 1. Start all containers in the background
docker compose up -d

# 2. View streaming logs from all services
docker compose logs -f

# 3. Verify health endpoints
curl http://localhost:8080/api/v1/health

# 4. Open the frontend application in your browser
# URL: http://localhost:3000
```

To tear down the containers and clean up:
```bash
docker compose down -v
```

---

## 3. Running Services Locally for Native Development

For day-to-day coding, running the services natively outside Docker provides instant hot-reloading and direct debugger attachment.

### Step 3.1: Start the PostgreSQL Database

If you have Docker installed, start only the database container:

```bash
docker compose up -d postgres
```

The database will be exposed on port `5432` with credentials:
* **Host**: `localhost`
* **Port**: `5432`
* **User**: `app_user`
* **Password**: `app_password`
* **Database**: `app_db`

---

### Step 3.2: Configure and Run the Go Backend

Navigate to the `backend/` directory:

```bash
cd backend

# Copy the environment template
cp .env.example .env

# Download Go dependencies
go mod download

# Run database migrations (or run backend server which automatically applies pending migrations on startup)
go run cmd/server/main.go
```

The backend server listens on `http://localhost:8080`.

Verify server liveness and readiness:
```bash
# Liveness probe (checks process responsiveness)
curl http://localhost:8080/api/v1/health/live

# Readiness probe (checks DB connectivity)
curl http://localhost:8080/api/v1/health/ready
```

---

### Step 3.3: Configure and Run the React Frontend

Open a second terminal and navigate to `frontend/`:

```bash
cd frontend

# Copy the environment template
cp .env.example .env

# Install dependencies cleanly
npm install

# Start Vite development server with Hot Module Replacement (HMR)
npm run dev
```

The frontend will run at `http://localhost:5173`. Requests to `/api/v1` are proxied or dispatched directly to `http://localhost:8080`.

---

## 4. Database Migrations and SQLC Generation

### Database Migrations
Migrations live in `backend/internal/database/migrations/`:
* `000001_init.up.sql`: Applies the initial schema and tables.
* `000001_init.down.sql`: Reverts the schema changes cleanly.

The backend bootstrap includes an embedded migration runner using Go's `embed` package, executing pending migrations safely inside an isolated transaction at application startup.

### SQLC Query Regeneration
When you add or update SQL queries in `backend/internal/database/queries/`:

```bash
cd backend
sqlc generate
```

SQLC will read `sqlc.yaml`, parse your queries against the migration schemas, and regenerate the typed Go files in `backend/internal/database/sqlc/`.

> **Note**: Never edit files in `backend/internal/database/sqlc/` manually. Always edit the source SQL in `queries/` or `migrations/` and run `sqlc generate`.

---

## 5. Testing and Quality Assurance

### Backend Quality Commands (from `backend/`)

```bash
# Format Go code
gofmt -s -w .

# Run static analysis
go vet ./...

# Run all unit and integration tests
go test -v ./...

# Run tests with the Go race detector enabled
go test -race ./...

# Build the production executable binary
go build -o bin/server cmd/server/main.go
```

### Frontend Quality Commands (from `frontend/`)

```bash
# Run ESLint to catch syntax and style issues
npm run lint

# Run TypeScript typechecker without emitting files
npm run typecheck

# Run Vitest test suites
npm run test

# Run Vitest with coverage report
npm run test:coverage

# Build optimized production bundle
npm run build
```

---

## 6. Environment Variables Reference

### Backend (`backend/.env`)

| Variable | Default Value | Description |
| :--- | :--- | :--- |
| `APP_ENV` | `development` | Environment mode (`development`, `staging`, `production`) |
| `HTTP_PORT` | `8080` | Port for the HTTP server to bind to |
| `DATABASE_URL` | `postgres://...` | Connection URI for PostgreSQL |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:5173` | Comma-separated allowed CORS origins |
| `READ_TIMEOUT_SECONDS` | `10` | Maximum duration for reading the entire request |
| `WRITE_TIMEOUT_SECONDS` | `10` | Maximum duration before timing out writes of the response |
| `IDLE_TIMEOUT_SECONDS` | `120` | Maximum amount of time to wait for the next request on keep-alive |
| `SHUTDOWN_TIMEOUT_SECONDS`| `15` | Maximum grace period for in-flight requests during shutdown |
| `DB_MAX_OPEN_CONNS` | `25` | Maximum number of open connections to the database pool |
| `DB_MIN_OPEN_CONNS` | `5` | Minimum number of idle connections maintained in pool |
| `DB_MAX_IDLE_TIME_MINUTES`| `15` | Maximum duration a connection may be idle before closing |
| `LOG_LEVEL` | `info` | Minimum log level (`debug`, `info`, `warn`, `error`) |
| `LOG_FORMAT` | `json` | Structured format (`json` or `text`) |

### Frontend (`frontend/.env`)

| Variable | Default Value | Description |
| :--- | :--- | :--- |
| `VITE_API_BASE_URL` | `http://localhost:8080/api/v1` | Public backend API URL |
| `VITE_APP_TITLE` | `App Bootstrap` | Display title for the web application |
