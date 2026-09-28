# Database Architecture & PostgreSQL Guidelines

## 1. Persistence Philosophy

This system pairs PostgreSQL with `pgx` (native PostgreSQL driver) and `sqlc` (compile-time type-safe Go code generator).

* **No Heavy ORMs**: We deliberately avoid ORMs (such as GORM or Hibernate-style mappers). ORMs obscure SQL execution, make N+1 queries easy to accidentally introduce, complicate connection pooling, and hinder query performance tuning.
* **SQL First**: All queries are authored by engineers in clean, standard PostgreSQL syntax in `backend/internal/database/queries/`.
* **Zero Run-Time Reflection for Queries**: SQLC parses the schema migrations and query files to generate compile-time checked Go structs and functions with zero runtime reflection overhead.

---

## 2. Directory Structure

```text
backend/internal/database/
├── migrations/
│   ├── 000001_init.up.sql       # Initial database schema definition
│   └── 000001_init.down.sql     # Safe rollback script
├── queries/
│   └── system.sql               # Hand-crafted SQL queries annotated for SQLC
├── sqlc/
│   ├── db.go                    # SQLC DBTX connection interface
│   ├── models.go                # Strongly typed Go structs representing tables
│   └── system.sql.go            # Generated Go query methods
├── pool.go                      # pgxpool lifecycle, health check, and connection pooling
└── migrate.go                   # Embedded migration runner executing at startup
```

---

## 3. Database Migration Workflow

1. **Plan & Review Schema Change**:
   Define the tables, constraints, types, and indexes needed.
2. **Create New Migration Files**:
   Create a new sequential migration pair under `backend/internal/database/migrations/`:
   * `<sequence>_<feature_name>.up.sql`
   * `<sequence>_<feature_name>.down.sql`
3. **Write Queries**:
   Author parameterized queries in `backend/internal/database/queries/<feature>.sql` with SQLC annotations (`-- name: GetSomething :one`, `-- name: ListItems :many`, etc.).
4. **Run SQLC Generation**:
   ```bash
   cd backend
   sqlc generate
   ```
5. **Verify Compilation & Run Tests**:
   ```bash
   go test -v ./...
   ```

---

## 4. Schema Design Principles

* **Database Constraints as First-Class Invariants**: Application validation is necessary for user experience, but database constraints provide absolute guarantees of data integrity.
  * Always use `NOT NULL` unless a column is genuinely optional.
  * Use foreign keys (`REFERENCES`) with explicit `ON DELETE` semantics (`RESTRICT`, `CASCADE`, or `SET NULL`).
  * Use `CHECK` constraints to enforce valid enum states, positive quantities, or non-empty strings.
  * Use `UNIQUE` constraints to enforce uniqueness at the database level to prevent race conditions.
* **Primary Key Strategy**:
  * For public-facing entities and distributed architectures: **UUIDv7** or **UUIDv4** (`UUID PRIMARY KEY DEFAULT gen_random_uuid()`) to prevent ID enumeration and ease multi-region replication.
  * For internal or append-only log tables: `BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY` for compact indexes and optimal B-tree page density.
* **Timestamps**:
  * Always use `TIMESTAMPTZ` (`TIMESTAMP WITH TIME ZONE`). Never store naive timestamps.
  * Standard audit columns: `created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP`, `updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP`.

---

## 5. Query Optimization & Best Practices

* **No Unbounded Queries**: Never write `SELECT ...` without a `LIMIT` clause on endpoints returning collections.
* **Explicit Column Projections**: Avoid `SELECT *`. Explicitly select only the columns needed by the calling context to reduce network I/O and memory allocations.
* **Covering Indexes**: Create indexes matching query filter patterns (`WHERE`, `ORDER BY`). For composite queries, order index columns from highest equality selectivity to range filters.
* **Avoid N+1 Queries**: Use SQL joins, subqueries, or aggregate functions (`json_agg` / `array_agg`) to retrieve related parent-child entities in a single round-trip.

---

## 6. Transactions & Concurrency

* **Transaction Scoping**: Transactions (`pgx.Tx`) must be created at the application service layer where business boundaries exist. Handlers and data access objects should never unilaterally start transactions without service coordination.
* **Atomic Unit of Work**:
  ```go
  tx, err := pool.Begin(ctx)
  if err != nil {
      return fmt.Errorf("begin tx: %w", err)
  }
  defer tx.Rollback(ctx) // Safe no-op if tx.Commit succeeded

  // Operations using SQLC WithTx(tx)
  qtx := queries.WithTx(tx)
  if err := qtx.InsertItem(ctx, ...); err != nil {
      return err
  }

  return tx.Commit(ctx)
  ```
* **Concurrency Handling**:
  * Do not rely on "read-then-write" in Go code to guard against concurrency conflicts.
  * Rely on PostgreSQL `UNIQUE` constraints with `ON CONFLICT DO UPDATE` or `ON CONFLICT DO NOTHING`.
  * For concurrent state modifications, use row-level locking (`SELECT ... FOR UPDATE`) or optimistic concurrency control via a `version INT NOT NULL DEFAULT 1` column.

---

## 7. Connection Pool Tuning (pgxpool)

Configured centrally in `backend/internal/database/pool.go`:

* `MaxConns`: Maximum open connections allowed in the pool (default: `25`).
* `MinConns`: Minimum idle connections preserved to absorb traffic spikes (default: `5`).
* `MaxConnIdleTime`: Duration before an unused idle connection is closed (default: `15m`).
* `MaxConnLifetime`: Maximum age of any connection before being recycled (default: `1h`).
* `HealthCheckPeriod`: Periodic interval for testing stale connections in the pool (default: `1m`).

---

## 8. Backup, Recovery & Disaster Planning

* **Automated Backups**: In production, configure automated daily base backups with continuous PostgreSQL Write-Ahead Logging (WAL) archiving (using tools like `pgBackRest` or managed cloud database automated backups).
* **Point-in-Time Recovery (PITR)**: Ensure WAL archives allow restoring the database to any millisecond within the backup retention window (typically 7 to 30 days).
* **Restore Drills**: Regularly verify database recovery by restoring backups to an isolated staging environment and executing health assertions.
