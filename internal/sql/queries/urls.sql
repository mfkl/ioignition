-- name: CreateSessionUrl :one
INSERT INTO urls (
  id, url, event_name, domain_session_id, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;
