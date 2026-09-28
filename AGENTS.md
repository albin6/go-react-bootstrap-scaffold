# Root AI Agent Guide: Go + React + PostgreSQL Production Scaffold

This repository is a production-grade, reusable full-stack bootstrap template intended to serve as the foundation for modern web applications.

## 1. Primary Architectural Invariant: Strict Separation

The repository is partitioned into two completely independent application trees:

```text
project-root/
├── frontend/     # Independent React, TypeScript, and Vite single-page application
├── backend/      # Independent Go REST API server, pgxpool, and SQLC persistence
├── docs/         # System architecture, standards, and engineering manuals
├── AGENTS.md     # Root agent guidelines (this file)
├── README.md     # Developer onboarding and project overview
└── .gitignore    # Repository-level ignore rules
```

### Absolute Boundary Rules
* **No Shared Code**: Do NOT create `shared/`, `common/`, `libs/`, or `packages/` directories at the project root.
* **No Cross-Boundary Imports**: The frontend must NEVER import backend Go files, SQL files, or templates. The backend must NEVER import frontend TypeScript, TSX, CSS, or JSON files.
* **No Root Package**: Never install npm packages or Go modules in the root directory.
* **Exclusive Transport**: Frontend communicates with the backend exclusively via HTTP requests over the network (`/api/v1/...`).

---

## 2. Technology Stack

* **Frontend**: React 18+, TypeScript (strict mode), Vite, React Router v6+, Axios / custom API client, Vitest, React Testing Library, ESLint, Prettier.
* **Backend**: Go (stable 1.23+), standard `net/http` with Go 1.22+ routing, PostgreSQL 15/16, `pgx/v5` connection pool (`pgxpool`), `sqlc` for type-safe Go SQL generation, embedded SQL migrations, standard `log/slog` structured logging, standard `testing` and `httptest`.
* **Database**: PostgreSQL with strict schemas, foreign keys, CHECK constraints, and transaction boundaries.

---

## 3. Directory Responsibilities

| Path | Purpose & Ownership |
| :--- | :--- |
| `frontend/src/features/` | Vertical domain slices (components, hooks, services, types). |
| `frontend/src/components/`| Reusable UI design system and layout components. |
| `frontend/src/lib/` | API client instance and third-party library wrappers. |
| `frontend/tests/` | Unit, component, and integration tests for frontend. |
| `backend/cmd/server/` | Main executable entrypoint (`main.go`). Kept minimal. |
| `backend/internal/config/`| Environment loader and validator. Only place `os.Getenv` is permitted. |
| `backend/internal/database/`| Database connection pool, migrations, queries, and SQLC code. |
| `backend/internal/http/` | HTTP middleware, routes, handlers, and JSON response envelopes. |
| `backend/internal/health/` | Production health check module (liveness & readiness). |
| `backend/internal/logger/` | Structured `log/slog` initialization and context correlation. |
| `docs/` | Comprehensive technical manuals for humans and agents. |

---

## 4. Coding Standards

### Go Rules
* **Idiomatic Go**: Prefer small concrete types and composition over class-like hierarchies.
* **Consumer-Defined Interfaces**: Define interfaces in the package where they are consumed, not where implemented. Keep them small (1-3 methods).
* **Context Propagation**: Always pass `ctx context.Context` as the first argument to functions performing I/O.
* **Error Handling**: Use `fmt.Errorf("...: %w", err)` for error wrapping. Never ignore errors.
* **No Unmanaged Goroutines**: Every goroutine must have an explicit lifecycle tied to context cancellation or sync groups.

### TypeScript / React Rules
* **Strict Mode**: Never use `any`. Use explicit interfaces for all props, states, and API contracts.
* **Feature Slices**: Place domain-specific UI and logic inside `src/features/<feature>/`.
* **API Client**: Never call `fetch()` directly in UI components. Use feature services.

---

## 5. Database & API Rules

* **SQL First**: Hand-craft SQL queries in `backend/internal/database/queries/`.
* **Never Edit SQLC Output**: Files in `backend/internal/database/sqlc/` are generated. Regenerate via `sqlc generate`.
* **Migrations**: All schema changes must be versioned in `backend/internal/database/migrations/`.
* **Standard JSON Envelope**: All API responses use the standardized `{ success: true, data: ... }` or `{ success: false, error: ... }` format.
* **Request Correlation**: Propagate `X-Request-ID` from ingress middleware to logs and response headers.

---

## 6. Agent Implementation Workflow

When executing any task:
1. Identify whether changes belong to `frontend/`, `backend/`, or both.
2. Read the corresponding `frontend/AGENTS.md` or `backend/AGENTS.md`.
3. If database changes are needed, write migration first -> write SQL query -> run `sqlc generate`.
4. Implement backend service and handler with tests.
5. Implement frontend types, service, hook, and UI component with tests.
6. Verify Go: `gofmt -s -w .`, `go vet ./...`, `go test -v ./...`.
7. Verify Frontend: `npm run lint`, `npm run typecheck`, `npm run test`, `npm run build`.
