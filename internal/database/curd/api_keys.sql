-- name: CreateAPIKey :one
INSERT INTO api_keys (
    user_id, provider_type, encrypted_key, key_hash, is_active
) VALUES (
    ?1, ?2, ?3, ?4, ?5
) RETURNING id, user_id, provider_type, encrypted_key, key_hash, is_active, created_at, updated_at;

-- name: GetAPIKey :one
SELECT id, user_id, provider_type, encrypted_key, key_hash, is_active, created_at, updated_at 
FROM api_keys
WHERE user_id = ?1 AND provider_type = ?2 AND is_active = TRUE
LIMIT 1;

-- name: GetAPIKeyByID :one
SELECT id, user_id, provider_type, encrypted_key, key_hash, is_active, created_at, updated_at 
FROM api_keys
WHERE id = ?1 LIMIT 1;

-- name: ListAPIKeysByUser :many
SELECT id, user_id, provider_type, encrypted_key, key_hash, is_active, created_at, updated_at 
FROM api_keys
WHERE user_id = ?1 AND is_active = TRUE
ORDER BY created_at DESC;

-- name: ListAPIKeysByProvider :many
SELECT id, user_id, provider_type, encrypted_key, key_hash, is_active, created_at, updated_at 
FROM api_keys
WHERE provider_type = ?1 AND is_active = TRUE
ORDER BY created_at DESC;

-- name: UpdateAPIKey :one
UPDATE api_keys 
SET encrypted_key = ?3, key_hash = ?4, updated_at = CURRENT_TIMESTAMP
WHERE user_id = ?1 AND provider_type = ?2 AND is_active = TRUE
RETURNING id, user_id, provider_type, encrypted_key, key_hash, is_active, created_at, updated_at;

-- name: DeactivateAPIKey :exec
UPDATE api_keys 
SET is_active = FALSE, updated_at = CURRENT_TIMESTAMP
WHERE user_id = ?1 AND provider_type = ?2;

-- name: DeleteAPIKey :exec
DELETE FROM api_keys
WHERE user_id = ?1 AND provider_type = ?2;

-- name: CountAPIKeysByUser :one
SELECT COUNT(*) FROM api_keys
WHERE user_id = ?1 AND is_active = TRUE;

-- name: CountAPIKeysByProvider :one
SELECT COUNT(*) FROM api_keys
WHERE provider_type = ?1 AND is_active = TRUE;

-- name: CheckAPIKeyExists :one
SELECT COUNT(*) FROM api_keys
WHERE user_id = ?1 AND provider_type = ?2 AND is_active = TRUE;