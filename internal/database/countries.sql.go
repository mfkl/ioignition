// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: countries.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createLocation = `-- name: CreateLocation :one
INSERT INTO countries (
  id, emoji, country_code, name, region, domain_session_id, domain_id, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING id, emoji, country_code, name, region, created_at, updated_at, domain_session_id, domain_id
`

type CreateLocationParams struct {
	ID              uuid.UUID
	Emoji           string
	CountryCode     string
	Name            string
	Region          string
	DomainSessionID uuid.UUID
	DomainID        uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (q *Queries) CreateLocation(ctx context.Context, arg CreateLocationParams) (Country, error) {
	row := q.db.QueryRowContext(ctx, createLocation,
		arg.ID,
		arg.Emoji,
		arg.CountryCode,
		arg.Name,
		arg.Region,
		arg.DomainSessionID,
		arg.DomainID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Country
	err := row.Scan(
		&i.ID,
		&i.Emoji,
		&i.CountryCode,
		&i.Name,
		&i.Region,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DomainSessionID,
		&i.DomainID,
	)
	return i, err
}

const getLocationCount = `-- name: GetLocationCount :many
SELECT name, emoji, COUNT(domain_session_id) AS location_count
FROM countries
WHERE domain_id = $1 AND created_at > $2
GROUP BY country_code, name, emoji
ORDER BY location_count DESC
`

type GetLocationCountParams struct {
	DomainID  uuid.UUID
	CreatedAt time.Time
}

type GetLocationCountRow struct {
	Name          string
	Emoji         string
	LocationCount int64
}

func (q *Queries) GetLocationCount(ctx context.Context, arg GetLocationCountParams) ([]GetLocationCountRow, error) {
	rows, err := q.db.QueryContext(ctx, getLocationCount, arg.DomainID, arg.CreatedAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLocationCountRow
	for rows.Next() {
		var i GetLocationCountRow
		if err := rows.Scan(&i.Name, &i.Emoji, &i.LocationCount); err != nil {
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
