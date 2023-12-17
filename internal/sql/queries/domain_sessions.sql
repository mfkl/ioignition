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
    SELECT generate_series(CURRENT_DATE - make_interval(days := @interval_start), CURRENT_DATE, make_interval(days := @step)) AS start_date
)
SELECT dr.start_date::TIMESTAMP, COUNT(DISTINCT ds.session_id) AS session_count
FROM date_ranges dr
LEFT JOIN domain_sessions ds ON ds.session_start_time >= dr.start_date AND ds.session_start_time < dr.start_date + make_interval(days := @step)
WHERE (ds.domain_id = $1 OR ds.domain_id IS NULL) -- return a null row, i.e. with 0 visits on days with no visitors
GROUP BY dr.start_date
ORDER BY dr.start_date;

-- name: GetRefererStats :many
SELECT referer, COUNT(referer) AS referer_count 
FROM domain_sessions
WHERE domain_id = $1 AND created_at > $2
GROUP BY referer
ORDER BY referer_count DESC;

-- name: GetPlatformStats :many
SELECT platform, COUNT(platform) AS platform_count 
FROM domain_sessions
WHERE domain_id = $1 AND created_at > $2
GROUP BY platform
ORDER BY platform_count DESC;

-- name: GetBrowserStats :many
SELECT browser, COUNT(browser) AS browser_count 
FROM domain_sessions
WHERE domain_id = $1 AND created_at > $2
GROUP BY browser
ORDER BY browser_count DESC;

-- name: GetCurrentlyActiveUsers :one
SELECT session_id, COUNT(DISTINCT session_id)
FROM domain_sessions
WHERE domain_id = $1 AND session_start_time > $2 AND session_end_time IS NULL
GROUP BY session_id;
