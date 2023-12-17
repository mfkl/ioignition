// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: domain_sessions.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createDomainSession = `-- name: CreateDomainSession :one
INSERT INTO domain_sessions (
  id, client_id, session_id, session_start_time, session_end_time, domain_id, referer, device_width, browser, platform, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING id, client_id, session_id, session_start_time, session_end_time, referer, device_width, browser, platform, created_at, updated_at, domain_id
`

type CreateDomainSessionParams struct {
	ID               uuid.UUID
	ClientID         string
	SessionID        string
	SessionStartTime time.Time
	SessionEndTime   sql.NullTime
	DomainID         uuid.UUID
	Referer          sql.NullString
	DeviceWidth      sql.NullInt32
	Browser          sql.NullString
	Platform         sql.NullString
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (q *Queries) CreateDomainSession(ctx context.Context, arg CreateDomainSessionParams) (DomainSession, error) {
	row := q.db.QueryRowContext(ctx, createDomainSession,
		arg.ID,
		arg.ClientID,
		arg.SessionID,
		arg.SessionStartTime,
		arg.SessionEndTime,
		arg.DomainID,
		arg.Referer,
		arg.DeviceWidth,
		arg.Browser,
		arg.Platform,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i DomainSession
	err := row.Scan(
		&i.ID,
		&i.ClientID,
		&i.SessionID,
		&i.SessionStartTime,
		&i.SessionEndTime,
		&i.Referer,
		&i.DeviceWidth,
		&i.Browser,
		&i.Platform,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DomainID,
	)
	return i, err
}

const getBrowserStats = `-- name: GetBrowserStats :many
SELECT browser, COUNT(browser) AS browser_count 
FROM domain_sessions
WHERE domain_id = $1 AND created_at > $2
GROUP BY browser
ORDER BY browser_count DESC
`

type GetBrowserStatsParams struct {
	DomainID  uuid.UUID
	CreatedAt time.Time
}

type GetBrowserStatsRow struct {
	Browser      sql.NullString
	BrowserCount int64
}

func (q *Queries) GetBrowserStats(ctx context.Context, arg GetBrowserStatsParams) ([]GetBrowserStatsRow, error) {
	rows, err := q.db.QueryContext(ctx, getBrowserStats, arg.DomainID, arg.CreatedAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBrowserStatsRow
	for rows.Next() {
		var i GetBrowserStatsRow
		if err := rows.Scan(&i.Browser, &i.BrowserCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCurrentlyActiveUsers = `-- name: GetCurrentlyActiveUsers :one
SELECT session_id, COUNT(DISTINCT session_id)
FROM domain_sessions
WHERE domain_id = $1 AND session_start_time > $2 AND session_end_time IS NULL
GROUP BY session_id
`

type GetCurrentlyActiveUsersParams struct {
	DomainID         uuid.UUID
	SessionStartTime time.Time
}

type GetCurrentlyActiveUsersRow struct {
	SessionID string
	Count     int64
}

func (q *Queries) GetCurrentlyActiveUsers(ctx context.Context, arg GetCurrentlyActiveUsersParams) (GetCurrentlyActiveUsersRow, error) {
	row := q.db.QueryRowContext(ctx, getCurrentlyActiveUsers, arg.DomainID, arg.SessionStartTime)
	var i GetCurrentlyActiveUsersRow
	err := row.Scan(&i.SessionID, &i.Count)
	return i, err
}

const getDomainSession = `-- name: GetDomainSession :one
SELECT id, client_id, session_id, session_start_time, session_end_time, referer, device_width, browser, platform, created_at, updated_at, domain_id FROM domain_sessions
WHERE client_id = $1 AND session_id = $2
`

type GetDomainSessionParams struct {
	ClientID  string
	SessionID string
}

func (q *Queries) GetDomainSession(ctx context.Context, arg GetDomainSessionParams) (DomainSession, error) {
	row := q.db.QueryRowContext(ctx, getDomainSession, arg.ClientID, arg.SessionID)
	var i DomainSession
	err := row.Scan(
		&i.ID,
		&i.ClientID,
		&i.SessionID,
		&i.SessionStartTime,
		&i.SessionEndTime,
		&i.Referer,
		&i.DeviceWidth,
		&i.Browser,
		&i.Platform,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DomainID,
	)
	return i, err
}

const getGraphStats = `-- name: GetGraphStats :many
WITH date_ranges AS (
    SELECT generate_series(CURRENT_DATE - make_interval(days := $3), CURRENT_DATE, make_interval(days := $2)) AS start_date
)
SELECT dr.start_date::TIMESTAMP, COUNT(DISTINCT ds.session_id) AS session_count
FROM date_ranges dr
LEFT JOIN domain_sessions ds ON ds.session_start_time >= dr.start_date AND ds.session_start_time < dr.start_date + make_interval(days := $2)
WHERE (ds.domain_id = $1 OR ds.domain_id IS NULL) -- return a null row, i.e. with 0 visits on days with no visitors
GROUP BY dr.start_date
ORDER BY dr.start_date
`

type GetGraphStatsParams struct {
	DomainID      uuid.UUID
	Step          int32
	IntervalStart int32
}

type GetGraphStatsRow struct {
	DrStartDate  time.Time
	SessionCount int64
}

func (q *Queries) GetGraphStats(ctx context.Context, arg GetGraphStatsParams) ([]GetGraphStatsRow, error) {
	rows, err := q.db.QueryContext(ctx, getGraphStats, arg.DomainID, arg.Step, arg.IntervalStart)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetGraphStatsRow
	for rows.Next() {
		var i GetGraphStatsRow
		if err := rows.Scan(&i.DrStartDate, &i.SessionCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPlatformStats = `-- name: GetPlatformStats :many
SELECT platform, COUNT(platform) AS platform_count 
FROM domain_sessions
WHERE domain_id = $1 AND created_at > $2
GROUP BY platform
ORDER BY platform_count DESC
`

type GetPlatformStatsParams struct {
	DomainID  uuid.UUID
	CreatedAt time.Time
}

type GetPlatformStatsRow struct {
	Platform      sql.NullString
	PlatformCount int64
}

func (q *Queries) GetPlatformStats(ctx context.Context, arg GetPlatformStatsParams) ([]GetPlatformStatsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPlatformStats, arg.DomainID, arg.CreatedAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPlatformStatsRow
	for rows.Next() {
		var i GetPlatformStatsRow
		if err := rows.Scan(&i.Platform, &i.PlatformCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRefererStats = `-- name: GetRefererStats :many
SELECT referer, COUNT(referer) AS referer_count 
FROM domain_sessions
WHERE domain_id = $1 AND created_at > $2
GROUP BY referer
ORDER BY referer_count DESC
`

type GetRefererStatsParams struct {
	DomainID  uuid.UUID
	CreatedAt time.Time
}

type GetRefererStatsRow struct {
	Referer      sql.NullString
	RefererCount int64
}

func (q *Queries) GetRefererStats(ctx context.Context, arg GetRefererStatsParams) ([]GetRefererStatsRow, error) {
	rows, err := q.db.QueryContext(ctx, getRefererStats, arg.DomainID, arg.CreatedAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRefererStatsRow
	for rows.Next() {
		var i GetRefererStatsRow
		if err := rows.Scan(&i.Referer, &i.RefererCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSessionStats = `-- name: GetSessionStats :one
SELECT 
  AVG(CASE WHEN session_start_time >= $2 THEN EXTRACT(EPOCH FROM (session_end_time - session_start_time)) ELSE NULL END) AS average_duration,
  COUNT(DISTINCT CASE WHEN session_start_time >= $2 THEN client_id ELSE NULL END) AS unique_visits,
  COUNT(CASE WHEN session_start_time >= $2 THEN 1 ELSE NULL END) AS total_visits,

  AVG(CASE WHEN session_start_time >= $3 AND session_start_time < $2 AND session_start_time IS NOT NULL THEN EXTRACT(EPOCH FROM (session_end_time - session_start_time)) ELSE 0.0 END) AS average_duration_prior,
  COUNT(DISTINCT CASE WHEN session_start_time >= $3 AND session_start_time < $2  THEN client_id ELSE NULL END) AS unique_visits_prior,
  COUNT(CASE WHEN session_start_time >= $3 AND session_start_time < $2  THEN 1 ELSE NULL END) AS total_visits_prior
FROM 
  domain_sessions
WHERE domain_id = $1
`

type GetSessionStatsParams struct {
	DomainID        uuid.UUID
	Interval        time.Time
	CompareInterval time.Time
}

type GetSessionStatsRow struct {
	AverageDuration      float64
	UniqueVisits         int64
	TotalVisits          int64
	AverageDurationPrior float64
	UniqueVisitsPrior    int64
	TotalVisitsPrior     int64
}

func (q *Queries) GetSessionStats(ctx context.Context, arg GetSessionStatsParams) (GetSessionStatsRow, error) {
	row := q.db.QueryRowContext(ctx, getSessionStats, arg.DomainID, arg.Interval, arg.CompareInterval)
	var i GetSessionStatsRow
	err := row.Scan(
		&i.AverageDuration,
		&i.UniqueVisits,
		&i.TotalVisits,
		&i.AverageDurationPrior,
		&i.UniqueVisitsPrior,
		&i.TotalVisitsPrior,
	)
	return i, err
}

const updateDomainSession = `-- name: UpdateDomainSession :exec
UPDATE domain_sessions 
  SET session_end_time = $2,
  updated_at = $3
WHERE id = $1
`

type UpdateDomainSessionParams struct {
	ID             uuid.UUID
	SessionEndTime sql.NullTime
	UpdatedAt      time.Time
}

func (q *Queries) UpdateDomainSession(ctx context.Context, arg UpdateDomainSessionParams) error {
	_, err := q.db.ExecContext(ctx, updateDomainSession, arg.ID, arg.SessionEndTime, arg.UpdatedAt)
	return err
}
