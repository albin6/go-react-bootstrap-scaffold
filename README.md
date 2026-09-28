# Go + React + PostgreSQL Production-Grade Bootstrap Scaffold

A clean, scalable, maintainable, production-ready full-stack bootstrap template engineered with strict architectural separation between an idiomatic Go backend and a modern React/TypeScript frontend.

```text
                    Browser
                       │
                       │ HTTPS / REST
                       ▼
              ┌─────────────────┐
              │ React Frontend  │
              │ TypeScript      │
              └────────┬────────┘
                       │
                       ▼
              ┌─────────────────┐
              │ Go HTTP API     │
              │ net/http        │
              └────────┬────────┘
                       │
                       ▼
              ┌─────────────────┐
              │ pgx + SQLC      │
              └────────┬────────┘
                       │
                       ▼
              ┌─────────────────┐
              │ PostgreSQL      │
              └─────────────────┘
```

---

## 1. Project Purpose

This scaffold serves as a reusable, production-ready starting foundation for full-stack client web applications. It deliberately omits domain-specific business features (like eCommerce, social feeds, or complex user management), providing instead a rock-solid, idiomatic foundation covering:

* Strict architectural boundaries and zero-dependency cross-coupling.
* Idiomatic Go HTTP backend using standard library `net/http` and Go 1.22+ routing.
* Native PostgreSQL persistence using `pgxpool` and compile-time type-safe code generation with `sqlc`.
* Feature-driven React frontend with TypeScript, Vite, React Router, Vitest, and testing libraries.
* Production readiness: structured logging (`log/slog`), request correlation (`X-Request-ID`), health probes (`/live` and `/ready`), graceful shutdown, security headers, and Docker containerization.

---

## 2. Technology Stack

### Frontend
* **Core**: React 18+, TypeScript (strict mode enabled)
* **Build & Dev Tooling**: Vite
* **Routing**: React Router v6+
* **Networking**: Centralized Axios HTTP client with typed envelopes
* **Linting & Formatting**: ESLint (Flat Config) & Prettier
* **Testing**: Vitest & React Testing Library

### Backend
* **Language & Runtime**: Go `1.23+` (compatible with stable `go1.27.1`)
* **Transport**: Standard library `net/http` with Go 1.22+ enhanced ServeMux routing
* **Database Driver**: `github.com/jackc/pgx/v5` and `pgxpool`
* **Query Layer**: `sqlc` for compile-safe Go query execution
* **Migrations**: Versioned SQL migrations (`.up.sql` and `.down.sql`)
* **Logging**: Standard library `log/slog` (structured JSON logging)
* **Testing**: Go standard `testing` and `httptest` packages

### Database & Infrastructure
* **Database**: PostgreSQL 16
* **Orchestration**: Docker & Docker Compose

---

## 3. Repository Structure

```text
project-root/
│
├── frontend/                     # Completely isolated React frontend application
│   ├── public/                   # Static web assets
│   ├── src/
│   │   ├── app/                  # Application root and shell
│   │   ├── assets/               # Bundled static media
│   │   ├── components/           # Generic reusable UI primitives & layout
│   │   ├── config/               # Frontend environment configuration
│   │   ├── features/             # Vertical domain slices (health, etc.)
│   │   ├── hooks/                # Reusable custom React hooks
│   │   ├── lib/                  # HTTP client and third-party library setup
│   │   ├── pages/                # Route view containers
│   │   ├── routes/               # Declarative client-side route configurations
│   │   ├── services/             # Base API client and shared network utils
│   │   ├── stores/               # Client-side state management
│   │   ├── types/                # Global frontend type definitions
│   │   ├── utils/                # Pure formatting and calculation utilities
│   │   └── main.tsx              # React entrypoint
│   ├── tests/                    # Vitest and RTL component tests
│   ├── .env.example              # Frontend environment variables template
│   ├── Dockerfile                # Multi-stage production build (Node -> Nginx)
│   ├── package.json              # Frontend dependencies and scripts
│   ├── tsconfig.json             # Strict TypeScript configuration
│   └── vite.config.ts            # Vite bundler and Vitest test config
│
├── backend/                      # Completely isolated Go backend service
│   ├── cmd/
│   │   └── server/
│   │       └── main.go           # Application bootstrap entrypoint
│   ├── internal/                 # Internal application packages (non-exportable)
│   │   ├── config/               # Environment configuration loader and validator
│   │   ├── database/             # PostgreSQL connection pool and SQLC layer
│   │   │   ├── migrations/       # Versioned SQL migrations (.up.sql / .down.sql)
│   │   │   ├── queries/          # Hand-crafted raw SQL queries for SQLC
│   │   │   ├── sqlc/             # Generated type-safe Go SQL code (DO NOT EDIT)
│   │   │   ├── pool.go           # pgxpool lifecycle and connection management
│   │   │   └── migrate.go        # Embedded SQL migration runner
│   │   ├── health/               # Health check module (liveness & readiness)
│   │   │   ├── handler.go        # HTTP handlers for /live, /ready, /health
│   │   │   ├── service.go        # Health status evaluation service
│   │   │   ├── types.go          # Domain health status types
│   │   │   └── handler_test.go   # httptest unit tests
│   │   ├── http/                 # Transport layer
│   │   │   ├── middleware/       # Correlation ID, logger, panic recovery, CORS, headers
│   │   │   ├── response/         # Standardized JSON response envelope helpers
│   │   │   └── routes/           # Mux registration and route wiring
│   │   └── logger/               # Structured slog initialization & context helpers
│   ├── tests/                    # Integration and helper tests
│   ├── sqlc.yaml                 # SQLC configuration
│   ├── go.mod                    # Go module dependencies
│   ├── go.sum                    # Go checksums
│   ├── .env.example              # Backend environment variables template
│   ├── Dockerfile                # Multi-stage minimal production container
│   └── README.md                 # Backend-specific documentation
│
├── docs/                         # In-depth architectural & operational manuals
│   ├── architecture.md           # System topology and architectural patterns
│   ├── development.md            # Local developer onboarding and commands
│   ├── coding-standards.md       # Go, TypeScript, and engineering standards
│   ├── api-conventions.md        # REST principles, envelopes, and error codes
│   ├── database.md               # PostgreSQL, pgx, SQLC, and migration guidelines
│   ├── security.md               # Security baseline and best practices
│   ├── observability.md          # Structured logging, correlation, and health probes
│   └── agent-guide.md            # Guidelines and rules for AI coding agents
│
├── docker-compose.yml            # Local development orchestration
├── AGENTS.md                     # Root AI Agent collaboration rules
├── README.md                     # Project overview (this file)
└── .gitignore                    # Repository-level ignore rules
```

---

## 4. Frontend / Backend Architectural Boundaries

1. **Physical Boundary**: `frontend/` and `backend/` are completely isolated trees.
2. **Zero Code Sharing**: Do not share source code between frontend and backend. No root-level package dependencies exist.
3. **Transport Exclusivity**: The frontend communicates with the Go backend solely over HTTP via REST APIs (`/api/v1/...`).
4. **Data Isolation**: The frontend has zero knowledge of SQL, PostgreSQL, pgx, or Go structs. The backend has zero knowledge of React, Vite, or client-side routing.

---

## 5. Prerequisites

* **Go**: Version `1.23+` (installed and accessible in PATH)
* **Node.js**: Version `20.x+` (LTS recommended)
* **npm**: Version `10.x+`
* **Docker & Docker Compose**: For containerized orchestration and PostgreSQL
* **SQLC** (Optional): `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

---

## 6. Local Setup and Quickstart

### Option A: Complete Stack via Docker Compose
```bash
# 1. Start all containers (PostgreSQL, Go backend, React frontend)
docker compose up -d

# 2. View logs
docker compose logs -f

# 3. Access endpoints
# Frontend: http://localhost:3000
# Backend Health: http://localhost:8080/api/v1/health
```

### Option B: Native Local Development

#### Step 1: Start PostgreSQL
```bash
docker compose up -d postgres
```

#### Step 2: Run Go Backend
```bash
cd backend
cp .env.example .env
go mod download
go run cmd/server/main.go
```
The Go server starts on `http://localhost:8080`.

#### Step 3: Run React Frontend
```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```
The React frontend starts on `http://localhost:5173`.

---

## 7. PostgreSQL Setup, Migrations & SQLC

### Migrations
Migrations are stored in `backend/internal/database/migrations/`:
* Applied automatically on application startup via an embedded Go runner.
* Written as plain, versioned SQL (`000001_init.up.sql`, `000001_init.down.sql`).

### SQLC Query Generation
When editing SQL queries in `backend/internal/database/queries/`:
```bash
cd backend
sqlc generate
```
This generates compile-time safe Go query code inside `backend/internal/database/sqlc/`.

---

## 8. Testing and Quality Checks

### Backend Quality Commands
```bash
cd backend
gofmt -s -w .       # Auto-format Go code
go vet ./...         # Static analysis
go test -v ./...     # Unit and HTTP handler tests
go test -race ./...  # Race condition detection
go build ./...       # Compile check
```

### Frontend Quality Commands
```bash
cd frontend
npm run lint         # ESLint checking
npm run typecheck    # TypeScript compiler check
npm run test         # Vitest unit & component test suite
npm run build        # Production bundle build
```

---

## 9. Observability & Health Probes

The backend exposes standardized operational endpoints under `/api/v1`:

* `GET /api/v1/health/live`: Fast liveness probe (verifies HTTP server process is responsive).
* `GET /api/v1/health/ready`: Dependency readiness probe (pings PostgreSQL database pool).
* `GET /api/v1/health`: Aggregate system health status with database status and uptime.

---

## 10. AI-Agent Context and Rules

If you are an AI coding agent working in this repository:
1. Refer to [AGENTS.md](file:///c:/My%20Files/New%20folder%20%283%29/AGENTS.md) and [docs/agent-guide.md](file:///c:/My%20Files/New%20folder%20%283%29/docs/agent-guide.md).
2. Never import backend code into the frontend or frontend code into the backend.
3. Never edit generated code in `backend/internal/database/sqlc/`.
4. Follow the step-by-step feature workflow outlined in [docs/agent-guide.md](file:///c:/My%20Files/New%20folder%20%283%29/docs/agent-guide.md).

---

## 11. Complete Documentation Index

* [System Architecture](file:///c:/My%20Files/New%20folder%20%283%29/docs/architecture.md)
* [Local Development Guide](file:///c:/My%20Files/New%20folder%20%283%29/docs/development.md)
* [Coding Standards](file:///c:/My%20Files/New%20folder%20%283%29/docs/coding-standards.md)
* [API Conventions](file:///c:/My%20Files/New%20folder%20%283%29/docs/api-conventions.md)
* [Database Guide](file:///c:/My%20Files/New%20folder%20%283%29/docs/database.md)
* [Security Baseline](file:///c:/My%20Files/New%20folder%20%283%29/docs/security.md)
* [Observability & Telemetry](file:///c:/My%20Files/New%20folder%20%283%29/docs/observability.md)
* [AI Agent Guide](file:///c:/My%20Files/New%20folder%20%283%29/docs/agent-guide.md)
