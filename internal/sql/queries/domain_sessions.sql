-- name: CreateDomainSession :one
INSERT INTO domain_sessions (
  id, client_id, session_id, session_start_time, session_end_time, domain_id, referer, device_width, browser, platform, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: GetDomainSession :one
SELECT * FROM domain_sessions
WHERE client_id = $1 AND session_id = $2;

-- name: UpdateDomainSession :exec
UPDATE domain_sessions 
  SET session_end_time = $2,
  updated_at = $3
WHERE id = $1;
