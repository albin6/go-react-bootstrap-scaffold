# AI Agent Engineering Guide

## 1. Prime Directive

You are working within a full-stack Go + React + PostgreSQL production scaffold designed for maintainability, idiomatic code, and strict architectural separation.

**Never violate the boundary rules:**
1. `frontend/` and `backend/` are strictly independent projects.
2. Never create shared code directories (`shared/`, `common/`, `libs/`).
3. Never import backend Go files or types into the frontend.
4. Never import frontend assets or code into the backend.
5. Frontend communicates with the backend exclusively via the versioned HTTP REST API (`/api/v1`).
6. Never manually edit generated files in `backend/internal/database/sqlc/`.
7. Never write business logic inside HTTP handlers or middleware.
8. Never expose internal database errors or secrets to HTTP responses.

---

## 2. Directory Responsibilities

| Directory | Ownership & Rules |
| :--- | :--- |
| `frontend/` | React, TypeScript, Vite, styles, tests, client-side routing. Zero knowledge of Go or PostgreSQL. |
| `backend/` | Go source code, pgx connection pool, SQLC config & queries, SQL migrations, server, tests. |
| `backend/internal/database/migrations/` | Versioned SQL schema changes (`.up.sql` and `.down.sql`). |
| `backend/internal/database/queries/` | Hand-authored, reviewable SQL queries for SQLC code generation. |
| `backend/internal/database/sqlc/` | Generated code. Read-only for humans and agents. Run `sqlc generate` to update. |
| `backend/internal/http/handlers/` | HTTP transport: deserialization, parameter extraction, calling services, serializing response. |
| `backend/internal/service/` or `modules/` | Pure business logic, transaction coordination. Accepts `context.Context`, returns domain errors. |
| `docs/` | System architecture, conventions, and engineering documentation. |

---

## 3. Step-by-Step Feature Implementation Workflow

When asked to implement a new feature:

```text
Step 1: Check boundaries
  → Determine the HTTP contract between frontend and backend.

Step 2: Database Schema & Migrations
  → Create new migration in backend/internal/database/migrations/
  → Write raw PostgreSQL schema with explicit constraints and indexes.

Step 3: Queries & SQLC
  → Add queries in backend/internal/database/queries/
  → Run `sqlc generate` from backend/ to regenerate type-safe Go code.

Step 4: Backend Service & Handler
  → Create domain service in backend/internal/modules/<feature>/
  → Create HTTP handler in backend/internal/http/handlers/ or module.
  → Register routes on the standard mux in backend/internal/http/routes/.

Step 5: Backend Verification
  → Run `gofmt -s -w .`
  → Run `go vet ./...`
  → Run `go test -v ./...` and `go test -race ./...`

Step 6: Frontend Feature Implementation
  → Add TypeScript types in frontend/src/features/<feature>/types/
  → Add API service call in frontend/src/features/<feature>/services/
  → Add custom React hook in frontend/src/features/<feature>/hooks/
  → Add UI components in frontend/src/features/<feature>/components/

Step 7: Frontend Verification
  → Run `npm run lint`
  → Run `npm run typecheck`
  → Run `npm run test`
  → Run `npm run build`
```

---

## 4. Prohibited Patterns

* **DO NOT** install npm packages in the repository root.
* **DO NOT** place Go dependencies in the repository root.
* **DO NOT** use `any` in TypeScript code.
* **DO NOT** create Java-style OOP hierarchies or factory boilerplate in Go.
* **DO NOT** use ORMs or generate dynamic raw SQL query strings.
* **DO NOT** use `os.Getenv()` outside `backend/internal/config/`.
* **DO NOT** commit `.env` files or API secrets.
