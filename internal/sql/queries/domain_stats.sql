-- name: CreateDomainStat :one
INSERT INTO domain_stats (
  id, url, referer, device_width, user_agent, domain_session_id, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;
