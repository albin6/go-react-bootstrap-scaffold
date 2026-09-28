-- name: GetSystemMetadataByKey :one
SELECT id, key, value, environment, status, created_at, updated_at
FROM system_metadata
WHERE key = $1
LIMIT 1;

-- name: UpsertSystemMetadata :one
INSERT INTO system_metadata (key, value, environment, status)
VALUES ($1, $2, $3, $4)
ON CONFLICT (key) DO UPDATE
SET value = EXCLUDED.value,
    environment = EXCLUDED.environment,
    status = EXCLUDED.status,
    updated_at = CURRENT_TIMESTAMP
RETURNING id, key, value, environment, status, created_at, updated_at;

-- name: CheckDatabaseHealth :one
SELECT 1 AS healthy;

-- name: ListSystemMetadata :many
SELECT id, key, value, environment, status, created_at, updated_at
FROM system_metadata
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
