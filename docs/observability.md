# Observability & Telemetry

## 1. Observability Architecture

Observability provides visibility into system behavior across all lifecycle phases through:
* **Structured Logging**: Using Go's standard `log/slog` package.
* **Distributed Tracing & Request Correlation**: Utilizing propagated `X-Request-ID` headers.
* **Health Checks & Probes**: Standardized liveness and readiness HTTP endpoints for container orchestrators (Kubernetes, AWS ECS, Nomad).

---

## 2. Structured Logging with `log/slog`

We use Go's standard library `log/slog` instead of heavyweight third-party logging packages.

### 2.1 Logger Configuration
The application initializes a global JSON or text structured logger on startup:
* **Development**: Pretty-printed text format or formatted JSON for readability.
* **Production**: Compact JSON output streaming to `stdout` for ingestion by log collectors (Fluentbit, Datadog, Grafana Loki, CloudWatch).

### 2.2 Standard Log Attributes
Every request log entry contains standard attributes:
```json
{
  "time": "2026-09-28T20:30:15.123Z",
  "level": "INFO",
  "msg": "HTTP request completed",
  "request_id": "c1f7a264-b03a-4efb-b892-e4d0d08a1001",
  "method": "GET",
  "path": "/api/v1/health",
  "status": 200,
  "duration_ms": 1.45,
  "remote_addr": "192.168.1.10:54320"
}
```

### 2.3 Contextual Logging Helper
Handlers and services extract the logger enriched with the request correlation ID:
```go
logger := logger.FromContext(ctx)
logger.Info("performing database health check", "target", "postgres")
```

### 2.4 Sensitive Data Redaction
The logger middleware and helper functions enforce automatic scrubbing of sensitive attributes:
* `authorization` headers
* passwords, secrets, and API keys
* sensitive query parameters (`?token=...`)
* connection strings containing credentials

---

## 3. Request Correlation & Distributed Tracing

1. **Extraction/Generation**: Inbound requests are inspected for `X-Request-ID`. If missing, a new UUID is generated.
2. **Context Binding**: The ID is stored in the Go `context.Context`.
3. **Egress Propagation**: Any outgoing HTTP or RPC calls made on behalf of the request forward the same `X-Request-ID`.
4. **Client Response**: The backend returns `X-Request-ID` in HTTP response headers, allowing frontend errors to reference the backend log entry.

---

## 4. Health & Readiness Probes

### 4.1 Liveness Probe (`GET /api/v1/health/live`)
* Checks if the Go HTTP server process is running, responsive, and able to schedule goroutines.
* Returns `200 OK` with `{ "status": "alive" }`.
* If this probe fails repeatedly, container orchestrators will restart the container.

### 4.2 Readiness Probe (`GET /api/v1/health/ready`)
* Checks whether the service is ready to accept user traffic by verifying upstream dependencies:
  * Tests PostgreSQL connectivity by executing a fast ping (`pool.Ping(ctx)`).
* Returns `200 OK` if all dependencies are healthy.
* Returns `503 Service Unavailable` if PostgreSQL is unreachable.
* Container orchestrators will pull the pod out of the load balancer pool until readiness succeeds.

### 4.3 Aggregate Health (`GET /api/v1/health`)
* Returns comprehensive system health metadata, including uptime, dependency checks, and service version.
