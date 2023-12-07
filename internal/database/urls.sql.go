// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: urls.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createSessionUrl = `-- name: CreateSessionUrl :one
INSERT INTO urls (
  id, url, event_name, domain_session_id, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING id, event_name, url, created_at, updated_at, domain_session_id
`

type CreateSessionUrlParams struct {
	ID              uuid.UUID
	Url             string
	EventName       string
	DomainSessionID uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (q *Queries) CreateSessionUrl(ctx context.Context, arg CreateSessionUrlParams) (Url, error) {
	row := q.db.QueryRowContext(ctx, createSessionUrl,
		arg.ID,
		arg.Url,
		arg.EventName,
		arg.DomainSessionID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Url
	err := row.Scan(
		&i.ID,
		&i.EventName,
		&i.Url,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DomainSessionID,
	)
	return i, err
}
