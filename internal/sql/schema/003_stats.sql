-- +goose Up
CREATE TABLE IF NOT EXISTS stats (
  id UUID PRIMARY KEY,
  visitor_id TEXT UNIQUE NOT NULL,
  session_start_time TIMESTAMP NOT NULL,
  session_end_time TIMESTAMP, -- can be null if user still active
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  url TEXT UNIQUE NOT NULL,
  referer TEXT,
  device_width INTEGER,
  user_agent TEXT,
  domain_id UUID NOT NULL,
  FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS stats;
