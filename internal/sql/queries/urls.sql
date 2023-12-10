-- name: CreateSessionUrl :one
INSERT INTO urls (
  id, url, event_name, domain_session_id, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetPageViewCount :one
SELECT 
  COUNT(CASE WHEN s.session_start_time >= @interval THEN 1 ELSE NULL END) AS total_page_views,
  COUNT(CASE WHEN s.session_start_time >= @compare_interval AND s.session_start_time < @interval  THEN 1 ELSE NULL END) AS total_page_views_prior
FROM urls u LEFT JOIN domain_sessions s
ON u.domain_session_id = s.id
WHERE s.domain_id = $1;
