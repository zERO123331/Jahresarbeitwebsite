CREATE INDEX IF NOT EXISTS updates_title_idx ON updates (lower(title), id);
CREATE INDEX IF NOT EXISTS updates_user_id_idx ON updates (user_id);
CREATE INDEX IF NOT EXISTS updates_created_at_idx ON updates (created_at);
CREATE INDEX IF NOT EXISTS updates_updated_at_idx ON updates (updated_at);