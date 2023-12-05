-- name: CreateDomainSession :one
INSERT INTO domain_sessions (
  id,client_id, event_id, session_start_time, session_end_time, domain_id, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetDomainSession :one
SELECT * FROM domain_sessions
WHERE client_id = $1 AND event_id = $2;

-- name: UpdateDomainSession :exec
UPDATE domain_sessions 
  SET session_end_time = $2,
  updated_at = $3
WHERE id = $1;
