CREATE INDEX IF NOT EXISTS tunes_title_idx ON tunes USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS tunes_keys_idx ON tunes USING GIN (keys);