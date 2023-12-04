-- name: CreateDomainSession :one
INSERT INTO domain_sessions (
  id, visitor_id, session_start_time, session_end_time, domain_id, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;
