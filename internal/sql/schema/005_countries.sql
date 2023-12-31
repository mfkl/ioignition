-- +goose Up
CREATE TABLE IF NOT EXISTS countries (
  id UUID PRIMARY KEY,
  emoji TEXT NOT NULL,
  country_code TEXT NOT NULL,
  name TEXT NOT NULL,
  region TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  domain_session_id UUID NOT NULL,
  domain_id UUID NOT NULL,
  FOREIGN KEY (domain_session_id) REFERENCES domain_sessions(id) ON DELETE CASCADE,
  FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS countries;
