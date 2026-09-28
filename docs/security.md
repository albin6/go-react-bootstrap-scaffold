# Security Architecture & Baseline

## 1. Core Principles

Security is treated as a foundational architectural property across both frontend and backend boundaries.

* **Defense in Depth**: Security controls operate at the network, reverse-proxy, HTTP transport, application service, and database levels.
* **Least Privilege**: Application database credentials have permissions strictly limited to runtime schemas. Containers execute as non-root users.
* **Zero Trust Transport**: All internal and external network communication must be secured with TLS.

---

## 2. HTTP Transport Security

### 2.1 Security Headers
The Go HTTP server automatically injects strict security headers on all responses via `internal/http/middleware/security.go`:

* `Strict-Transport-Security: max-age=63072000; includeSubDomains; preload` (enforces HTTPS)
* `X-Content-Type-Options: nosniff` (prevents MIME type sniffing)
* `X-Frame-Options: DENY` (mitigates clickjacking)
* `X-XSS-Protection: 1; mode=block` (legacy XSS filtering)
* `Referrer-Policy: strict-origin-when-cross-origin` (prevents referrer leak)
* `Content-Security-Policy: default-src 'self'; frame-ancestors 'none';` (restricts content sources)

### 2.2 Cross-Origin Resource Sharing (CORS)
* CORS is strictly enforced via `internal/http/middleware/cors.go`.
* Only explicitly configured origins in `CORS_ALLOWED_ORIGINS` are granted access.
* Wildcard origins (`*`) are disallowed in production environments.
* Preflight `OPTIONS` requests are handled safely with appropriate `Access-Control-Max-Age` caching.

### 2.3 Request Limiting & Denial-of-Service Defense
* **Request Body Limits**: The middleware enforces a default maximum request body size (e.g., 1MB) via `http.MaxBytesReader` to protect against memory exhaustion attacks.
* **Server Timeouts**: Timeouts are explicitly configured on `http.Server`:
  * `ReadTimeout`: 10 seconds (guards against slowloris attacks)
  * `ReadHeaderTimeout`: 5 seconds
  * `WriteTimeout`: 10 seconds
  * `IdleTimeout`: 120 seconds

---

## 3. Data Protection & SQL Injection Prevention

* **Parameterized SQL via SQLC**: All SQL statements are pre-compiled and parameterized. String interpolation or raw query concatenation is strictly forbidden.
* **No Secret Leakage**:
  * Passwords, tokens, database connection strings, and API keys are scrubbed before logging.
  * Internal error messages and SQL diagnostics are logged server-side only; clients receive sanitized domain error codes.
* **Environment Secret Management**: Secrets are supplied via environment variables or secret managers (e.g., AWS Secrets Manager, HashiCorp Vault, Kubernetes Secrets). They are never committed to version control.

---

## 4. Frontend Security Safeguards

* **Vite Environment Boundaries**: Only variables prefixed with `VITE_` are bundled into the client build. No private keys, database credentials, or secret signing keys may ever be placed in `frontend/.env`.
* **Cross-Site Scripting (XSS) Prevention**: React inherently escapes string interpolation in JSX. The use of `dangerouslySetInnerHTML` is prohibited unless audited and sanitized with DOMPurify.
* **Token Storage Recommendations**: When authentication is integrated, store session tokens in `HttpOnly`, `Secure`, `SameSite=Strict` cookies to guard against client-side XSS theft.

---

## 5. Future Authentication & Authorization Boundary

When implementing authentication in this scaffold:
1. **Transport**: Use `HttpOnly` session cookies or short-lived Bearer tokens.
2. **Context**: Validate tokens in an authentication middleware (`internal/http/middleware/auth.go`) and inject an authenticated `UserContext` into the request `ctx`.
3. **Authorization**: Enforce role-based access control (RBAC) at the domain service boundary, ensuring authorization checks cannot be bypassed by alternate transports.
