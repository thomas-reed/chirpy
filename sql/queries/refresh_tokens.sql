-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, created_at, updated_at, expires_at, revoked_at)
VALUES (
  $1,
  $2,
  NOW(),
  NOW(),
  NOW() + INTERVAL '60 days',
  NULL
)
RETURNING user_id, created_at, updated_at, expires_at;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1
RETURNING token, user_id, updated_at, revoked_at;