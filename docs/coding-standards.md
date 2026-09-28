# Coding Standards & Engineering Guidelines

## 1. Go Backend Standards

### 1.1 Idiomatic Go Design (No Java / Node.js Parallels)
* **No Class Hierarchies**: Go is not an object-oriented language with inheritance. Do not create inheritance trees, abstract base structs, or factory-of-factories.
* **Avoid Premature Interfaces**: Do not create interfaces in the same package where concrete types are defined unless multiple implementations exist today. Follow the Go proverb: *"Accept interfaces, return structs."*
* **Consumer-Defined Interfaces**: Interfaces must be defined where they are consumed, not where they are implemented. Keep interfaces minimal (1 to 3 methods maximum).
* **Package Cohesion**: Each package must have a singular purpose. Never create grab-bag packages named `utils`, `helpers`, `common`, `misc`, or `manager`.
* **Explicit Dependency Injection**: Inject dependencies as explicit fields in constructors (e.g., `NewService(...)`). Never use global mutable singletons or runtime reflection-based container libraries.

### 1.2 Error Handling
* Errors are values. Always handle errors explicitly at the call site:
  ```go
  res, err := s.db.GetUser(ctx, id)
  if err != nil {
      return nil, fmt.Errorf("get user %s: %w", id, err)
  }
  ```
* Use `%w` with `fmt.Errorf` to wrap errors and preserve error chains for `errors.Is` and `errors.As`.
* Never discard errors with `_ = ...` unless deliberately documented as safe.
* Never expose raw SQL errors, stack traces, or internal file paths to HTTP clients. Translate errors into standardized domain errors before generating HTTP responses.

### 1.3 Context Propagation
* The first parameter of any function performing I/O, database access, or external requests must be `ctx context.Context`.
* Respect cancellation: If `ctx.Done()` fires, terminate work immediately and release resources.
* Never instantiate `context.Background()` inside request handlers or service methods—propagate the caller's request context.
* Do not use `context.WithValue` to pass optional function arguments or domain entities. Use it strictly for request-scoped operational metadata (e.g., `request_id`).

### 1.4 Concurrency & Goroutine Safety
* Every goroutine spawned must have an explicit lifecycle tied to context cancellation or a `sync.WaitGroup`.
* Never spawn unmanaged background goroutines inside HTTP handlers (`go doSomething()`). If an HTTP request finishes or times out, an unmanaged goroutine will leak or attempt to write to an aborted connection.
* Always run `go test -race ./...` to detect race conditions before committing code.

---

## 2. SOLID, DRY, and KISS in Go

### Single Responsibility Principle (SRP)
* A package, struct, or function should do one thing well. Handlers handle HTTP transport; services handle domain logic; SQLC handles SQL execution.

### Open/Closed Principle (OCP)
* Extend functionality using composition, struct embedding, and functional options rather than modifying established core structures.

### Liskov Substitution & Interface Segregation (LSP & ISP)
* Design tiny, targeted interfaces (e.g., `io.Reader`, `io.Writer`). Consumers define the exact interface subset they need.

### Dependency Inversion Principle (DIP)
* High-level domain logic should not be coupled directly to concrete third-party network clients or driver implementations; rely on consumer-defined interfaces when substitution is required.

### DRY (Don't Repeat Yourself)
* Avoid duplicating complex domain logic or business rules. However, do not eliminate incidental duplication if doing so couples unrelated domains or degrades readability.

### KISS (Keep It Simple, Stupid)
* Prefer clear, straightforward, readable code over clever one-liners or premature abstractions.

---

## 3. TypeScript & React Frontend Standards

### 3.1 Strict TypeScript
* `noImplicitAny`, `strictNullChecks`, and `strict` mode are permanently enabled in `tsconfig.json`.
* The `any` type is strictly prohibited. Use `unknown`, generic parameters, or explicit interfaces/type unions.
* All API response payloads must be typed using explicit TypeScript interfaces matching backend contracts.

### 3.2 Component Architecture & State
* Use pure functional components with explicit props interfaces (`interface Props { ... }`).
* Group UI into:
  * `components/ui/`: Atomic design system components (Button, Card, Badge, Modal).
  * `components/layout/`: Page scaffolding (Navbar, Container, Footer).
  * `features/<name>/components/`: Domain-specific components bound to feature data.
* Keep state as local as possible. Use React state hooks (`useState`, `useReducer`) for UI state and custom feature hooks for data-fetching state.
* Centralize global state only when multiple detached components require access to identical data.

### 3.3 HTTP & API Consumption
* Never call `fetch()` or `axios.get()` directly inside React UI components.
* Route all HTTP calls through dedicated feature services (`src/features/<feature>/services/`) which utilize the centralized HTTP client (`src/lib/api-client.ts`).
* Normalize all backend errors into predictable frontend error models.

### 3.4 Formatting & Linting
* Code is automatically formatted with Prettier and linted with ESLint (using Flat Config).
* Code must pass `npm run lint` and `npm run typecheck` without warnings or errors.
