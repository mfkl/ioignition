-- +goose Up
CREATE TABLE IF NOT EXISTS domain_stats (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  url TEXT UNIQUE NOT NULL,
  referer TEXT,
  device_width INTEGER,
  user_agent TEXT,
  domain_session_id UUID NOT NULL,
  FOREIGN KEY (domain_session_id) REFERENCES domain_sessions(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS domain_stats;
