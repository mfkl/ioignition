-- name: CreateLocation :one
INSERT INTO countries (
  id, emoji, country_code, name, region, domain_session_id, domain_id, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetLocationCount :many
SELECT name, emoji, COUNT(domain_session_id) AS location_count
FROM countries
WHERE domain_id = $1 AND created_at > $2
GROUP BY country_code, name, emoji
ORDER BY location_count DESC;
