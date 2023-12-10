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

-- name: GetSessionStats :one
SELECT 
  AVG(CASE WHEN session_start_time >= @interval THEN EXTRACT(EPOCH FROM (session_end_time - session_start_time)) ELSE NULL END) AS average_duration,
  COUNT(DISTINCT CASE WHEN session_start_time >= @interval THEN client_id ELSE NULL END) AS unique_visits,
  COUNT(CASE WHEN session_start_time >= @interval THEN 1 ELSE NULL END) AS total_visits,

  AVG(CASE WHEN session_start_time >= @compare_interval AND session_start_time < @interval AND session_start_time IS NOT NULL THEN EXTRACT(EPOCH FROM (session_end_time - session_start_time)) ELSE 0.0 END) AS average_duration_prior,
  COUNT(DISTINCT CASE WHEN session_start_time >= @compare_interval AND session_start_time < @interval  THEN client_id ELSE NULL END) AS unique_visits_prior,
  COUNT(CASE WHEN session_start_time >= @compare_interval AND session_start_time < @interval  THEN 1 ELSE NULL END) AS total_visits_prior
FROM 
  domain_sessions
WHERE domain_id = $1;

-- name: GetGraphStats :many
WITH date_ranges AS (
    SELECT generate_series(@interval_start::TIMESTAMP, CURRENT_DATE, make_interval(days := @step)) AS start_date
)
SELECT dr.start_date::TIMESTAMP, COUNT(DISTINCT ds.session_id) AS session_count
FROM date_ranges dr
LEFT JOIN domain_sessions ds ON ds.session_start_time >= dr.start_date AND ds.session_start_time < dr.start_date + make_interval(days := @step)
WHERE (ds.domain_id = $1 OR ds.domain_id IS NULL) AND dr.start_date <= CURRENT_DATE
GROUP BY dr.start_date
ORDER BY dr.start_date;
