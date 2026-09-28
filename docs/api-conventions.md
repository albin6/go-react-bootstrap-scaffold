# API Conventions & REST Contract

## 1. RESTful Principles & URI Design

All HTTP endpoints are versioned under the `/api/v1` namespace.

* **Resource Names**: Plural nouns in lowercase kebab-case (e.g., `/api/v1/health`, `/api/v1/system-metrics`).
* **HTTP Methods**:
  * `GET`: Safe, idempotent read operations. Must not mutate state.
  * `POST`: Create a new resource or invoke an operation.
  * `PUT`: Completely replace an existing resource.
  * `PATCH`: Partially update an existing resource.
  * `DELETE`: Remove a resource. Idempotent.

---

## 2. Standardized JSON Envelopes

Every JSON response from the Go API conforms to a predictable, strongly typed envelope structure.

### 2.1 Success Response (`200 OK`, `201 Created`, etc.)

```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "timestamp": "2026-09-28T20:30:00Z",
    "version": "1.0.0",
    "checks": {
      "database": "up"
    }
  }
}
```

For list endpoints requiring pagination:
```json
{
  "success": true,
  "data": {
    "items": [...],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total_items": 142,
      "total_pages": 8,
      "has_next": true
    }
  }
}
```

### 2.2 Error Response (`4xx`, `5xx`)

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "The payload contains invalid field values.",
    "details": [
      {
        "field": "email",
        "issue": "Must be a valid email address."
      }
    ]
  }
}
```

---

## 3. Standard HTTP Status Codes & Error Classifications

| HTTP Status | Error Code | Description |
| :--- | :--- | :--- |
| `200 OK` | N/A | Successful read or idempotent update. |
| `201 Created` | N/A | Resource successfully created. |
| `204 No Content`| N/A | Successful mutation with no response body. |
| `400 Bad Request` | `BAD_REQUEST` | Malformed JSON, unparseable input. |
| `401 Unauthorized` | `UNAUTHORIZED` | Missing or invalid authentication token. |
| `403 Forbidden` | `FORBIDDEN` | Authenticated caller lacks permission. |
| `404 Not Found` | `NOT_FOUND` | The requested resource or route does not exist. |
| `409 Conflict` | `CONFLICT` | Resource state conflict (e.g., unique key violation). |
| `422 Unprocessable`| `VALIDATION_FAILED`| Structurally valid JSON failing domain validation. |
| `429 Too Many Req`| `RATE_LIMITED` | Client has exceeded quota/rate limits. |
| `500 Internal Error`| `INTERNAL_ERROR` | Unexpected server failure. Diagnostic logged server-side. |
| `503 Unavailable`| `SERVICE_UNAVAILABLE`| Downstream dependency (e.g. database) is failing healthcheck. |

---

## 4. Request Correlation & Tracing

* Every inbound request receives an `X-Request-ID` header.
* If provided by a trusted upstream reverse proxy (e.g., Cloudflare, Envoy), the server validates and adopts it.
* If absent or invalid, the backend's correlation middleware generates a standard UUID.
* The `X-Request-ID` is:
  1. Attached to the request `context.Context`.
  2. Included in every structured log line via `slog`.
  3. Returned in the HTTP response header `X-Request-ID`.

---

## 5. Pagination, Sorting, and Filtering Standards

For future list endpoints, adhere to standard query parameter conventions:

* **Pagination**:
  * Offset-based: `?page=1&limit=20` (defaults: `page=1`, `limit=20`, max `limit=100`).
  * Cursor-based: `?cursor=eyJpZCI6MTAwfQ&limit=20` (recommended for high-volume append-only datasets).
* **Sorting**: `?sort=created_at&order=desc` (always restrict sortable fields to an explicit allowlist).
* **Filtering**: Standard equality checks: `?status=active`. Never execute unfiltered queries against PostgreSQL.
