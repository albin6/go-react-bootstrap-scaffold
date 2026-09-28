# Backend AI Agent Guide

## 1. Absolute Directives

* **Backend must never import frontend source code.**
* **Database access belongs to the backend.**
* **PostgreSQL is accessed through the database layer.**
* **SQLC generated code must not be manually edited.**
* **Schema changes require migrations.**
* **Handlers must not contain business logic.**
* **Services must not depend on HTTP request/response objects.**
* **Do not create interfaces unless they provide a concrete architectural or testing benefit.**

---

## 2. Architectural Boundaries & Principles

* **Package Boundaries**: Keep packages focused. Handlers live in `internal/http/handlers/` or feature modules (like `internal/health/`), services in domain modules, and database access in `internal/database/`.
* **No Java OOP Boilerplate**: Do not create factory structs, generic repositories, or service locators. Accept interfaces where consumed; return concrete structs.
* **Context Propagation**: Always pass `ctx context.Context` through handler -> service -> database. Never create `context.Background()` within request handlers.
* **Error Wrapping**: Wrap errors with `fmt.Errorf("operation: %w", err)`. Never return raw database errors or stack traces in HTTP responses.
* **Concurrency**: Always clean up goroutines with context cancellation. Always verify code with `go test -race ./...`.
* **Configuration**: Only read environment variables in `internal/config/config.go`. Never call `os.Getenv` in domain code.

---

## 3. Database Workflow for Agents

1. Add table definition or schema alter statement in `internal/database/migrations/<seq>_<name>.up.sql`.
2. Add reversal script in `internal/database/migrations/<seq>_<name>.down.sql`.
3. Add parameterized SQL in `internal/database/queries/<domain>.sql`.
4. Run `sqlc generate`.
5. Write domain service calling the generated SQLC methods within transaction boundaries where appropriate.
6. Write unit tests using `httptest` or integration tests.

---

## 4. Verification Checklist

Before finishing any backend modification:
* [ ] `gofmt -s -w .` passes without changes.
* [ ] `go vet ./...` reports zero issues.
* [ ] `go test -v ./...` passes.
* [ ] `go test -race ./...` passes with zero race conditions.
* [ ] `go build ./...` compiles cleanly.
