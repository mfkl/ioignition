-- name: CreateDomain :one
INSERT INTO domains (
  id, url, user_id, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: ListDomains :many
SELECT * FROM domains 
WHERE user_id = $1
LIMIT 10;
