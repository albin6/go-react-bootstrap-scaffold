# System Architecture

## 1. Architectural Philosophy

This full-stack scaffold is built upon strict architectural boundaries, explicit control flow, and idiomatic separation of concerns.

The core principles driving this architecture are:

* **Strict Frontend/Backend Separation**: The frontend and backend live in completely independent directory trees (`frontend/` and `backend/`). They share zero source code, dependencies, or runtime assumptions. All communication happens strictly through standard versioned HTTP REST APIs.
* **Idiomatic Go Backend**: We adhere strictly to idiomatic Go. We do not transplant Java Spring or Node.js/Express patterns into Go. We avoid unnecessary interface abstractions, service locators, dependency injection frameworks, and monolithic "clean architecture" layers that add indirection without value. We use small concrete types, accept interfaces only at consumption boundaries, rely on `context.Context` for request lifecycle propagation, and leverage the standard `net/http` package.
* **Explicit, Type-Safe Data Persistence**: We avoid heavy Object-Relational Mappers (ORMs). The persistence layer relies on PostgreSQL with `pgx` (for performant, connection-pooled native database interactions) and `sqlc` (for generating type-safe, compiled Go models and query executors directly from reviewed raw SQL).
* **Feature-Oriented React Frontend**: The frontend is built with React, TypeScript, and Vite. Code is organized into self-contained feature slices under `src/features/`, isolating UI components, client-side hooks, domain services, and types from unrelated concerns.
* **Statelessness and Horizontal Scalability**: The backend retains zero in-memory user session state, allowing any number of backend instances to run behind a round-robin load balancer.

---

## 2. High-Level System Topology

```text
                           Client Browser
                                │
                                │ HTTPS / TLS
                                ▼
                       ┌─────────────────┐
                       │ React Frontend  │
                       │ (Vite / SPA)    │
                       └────────┬────────┘
                                │
                                │ HTTP / JSON REST
                                ▼
                       ┌─────────────────┐
                       │  Load Balancer  │
                       │ (Reverse Proxy) │
                       └────────┬────────┘
                                │
                ┌───────────────┼───────────────┐
                ▼               ▼               ▼
         ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
         │ Go Backend  │ │ Go Backend  │ │ Go Backend  │
         │ Instance 1  │ │ Instance 2  │ │ Instance N  │
         └──────┬──────┘ └──────┬──────┘ └──────┬──────┘
                │               │               │
                └───────────────┼───────────────┘
                                │
                                │ pgxpool
                                ▼
                       ┌─────────────────┐
                       │   PostgreSQL    │
                       │  Primary DB     │
                       └─────────────────┘
```

---

## 3. Directory Layout & Boundaries

```text
project-root/
│
├── frontend/                     # Completely isolated React/TypeScript app
│   ├── public/                   # Static assets served by web server
│   ├── src/
│   │   ├── app/                  # Application root & shell
│   │   ├── assets/               # Bundled images, icons, global styles
│   │   ├── components/           # Reusable generic UI widgets
│   │   │   ├── common/           # Shared cross-cutting components (ErrorBoundary, etc.)
│   │   │   ├── layout/           # App layout frames, headers, footers
│   │   │   └── ui/               # Core design system primitives (Button, Card, Input)
│   │   ├── config/               # Frontend runtime configuration (env parsing)
│   │   ├── features/             # Domain-specific vertical slices
│   │   │   └── health/           # System health / bootstrap validation feature
│   │   ├── hooks/                # Cross-cutting custom React hooks
│   │   ├── lib/                  # Configured HTTP client & third-party adapters
│   │   ├── pages/                # Route view containers
│   │   ├── routes/               # Declarative client-side routing definitions
│   │   ├── services/             # Base API client and shared network utilities
│   │   ├── stores/               # Client-side state management
│   │   ├── types/                # Global frontend type definitions
│   │   ├── utils/                # Pure formatting and calculation helpers
│   │   └── main.tsx              # Application entrypoint
│   ├── tests/                    # Vitest and React Testing Library setup/suites
│   ├── .env.example              # Template for public Vite environment variables
│   ├── Dockerfile                # Multi-stage production build (Node -> Nginx Alpine)
│   ├── package.json              # Frontend dependencies and scripts
│   └── vite.config.ts            # Vite bundler & test configuration
│
├── backend/                      # Completely isolated Go backend service
│   ├── cmd/
│   │   └── server/
│   │       └── main.go           # Application executable bootstrap entrypoint
│   ├── internal/                 # Private application implementation
│   │   ├── config/               # Environment configuration loader and validator
│   │   ├── database/             # PostgreSQL connection pool and SQLC layer
│   │   │   ├── migrations/       # Versioned SQL migrations (up/down)
│   │   │   ├── queries/          # Hand-authored SQL queries for SQLC
│   │   │   ├── sqlc/             # Generated type-safe Go SQL code (read-only)
│   │   │   └── pool.go           # pgxpool lifecycle and health check
│   │   ├── health/               # Health check domain module (liveness/readiness)
│   │   │   ├── handler.go        # HTTP endpoints for /live, /ready, /health
│   │   │   ├── service.go        # Health status business evaluation
│   │   │   └── types.go          # Domain types and statuses
│   │   ├── http/                 # Transport layer
│   │   │   ├── middleware/       # Correlation ID, logger, panic recovery, CORS, headers
│   │   │   ├── response/         # Standardized JSON response envelopes & error formatting
│   │   │   └── routes/           # Mux registration and route mapping
│   │   └── logger/               # Structured slog initialization & context helpers
│   ├── tests/                    # Integration and end-to-end testing utilities
│   ├── sqlc.yaml                 # SQLC code generation configuration
│   ├── go.mod                    # Go module dependencies
│   ├── go.sum                    # Checksums for Go dependencies
│   ├── .env.example              # Backend environment template
│   └── Dockerfile                # Multi-stage minimal scratch/alpine container
│
├── docs/                         # Comprehensive engineering documentation
├── docker-compose.yml            # Local development orchestration
├── AGENTS.md                     # AI Agent collaboration instructions
├── README.md                     # Project overview and developer onboarding
└── .gitignore                    # Root repository ignore rules
```

---

## 4. Backend Request Flow

Requests flow through clear, single-responsibility boundaries:

```text
HTTP Request
      │
      ▼
HTTP Middleware Chain
  1. Request ID (Correlation tracking)
  2. Panic Recovery (Fault isolation)
  3. Structured Request Logger (slog)
  4. Security Headers (HSTS, CSP, X-Frame-Options, etc.)
  5. CORS Policy
  6. Request Timeout & Body Limit
      │
      ▼
net/http.ServeMux (Go 1.22+ Method & Path Pattern Matcher)
      │
      ▼
Handler (Transport concerns: decode JSON, validate parameters, invoke service)
      │
      ▼
Domain Service (Application business logic, transaction boundaries)
      │
      ▼
SQLC / pgxpool (Type-safe SQL queries against PostgreSQL)
      │
      ▼
Handler (Encapsulate outcome in standardized JSON response envelope)
      │
      ▼
HTTP Response
```

### Handler Constraints
* Handlers never construct raw SQL queries or touch `pgxpool.Pool` directly.
* Handlers never contain domain validation or business logic.
* Handlers translate transport parameters into Go types, call the domain service, and encode the response using `internal/http/response`.

### Service Constraints
* Services are decoupled from HTTP transport. They never accept `*http.Request` or return `http.ResponseWriter`.
* All service methods accept a `context.Context` as their first argument and pass it down to all database or external network calls.
* Services return standard Go errors. They do not assume an HTTP status code; the handler maps domain error classifications to HTTP status codes.

---

## 5. PostgreSQL Persistence Strategy

* **PostgreSQL + pgx**: We use the native `github.com/jackc/pgx/v5` driver and connection pool (`pgxpool.Pool`). Pgx supports native PostgreSQL binary protocol formats, prepared statements, and advanced data types.
* **SQLC Code Generation**: SQL queries are written in plain, reviewable `.sql` files in `backend/internal/database/queries/`. SQLC parses these queries against the schema migrations and produces compile-safe Go structs and execution methods.
* **Migrations**: All schema evolution is managed via sequential, bidirectional SQL migration scripts (`.up.sql` and `.down.sql`). Destructive schema synchronization or automatic ORM schema migrations in production are strictly forbidden.
