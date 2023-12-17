-- +goose Up
CREATE TABLE IF NOT EXISTS domain_sessions (
  id UUID PRIMARY KEY,
  client_id TEXT NOT NULL,
  -- sessionId helps identify a single session, comes from client
  session_id TEXT UNIQUE NOT NULL,
  session_start_time TIMESTAMP NOT NULL,
  session_end_time TIMESTAMP, -- can be null if user still active
  referer TEXT,
  device_width INTEGER,
  browser TEXT,
  platform TEXT,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  domain_id UUID NOT NULL,
  FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS domain_sessions;
