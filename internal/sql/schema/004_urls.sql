-- +goose Up
CREATE TABLE IF NOT EXISTS urls (
  id UUID PRIMARY KEY,
  url TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  domain_session_id UUID NOT NULL,
  FOREIGN KEY (domain_session_id) REFERENCES domain_sessions(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS urls;
