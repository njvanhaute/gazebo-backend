CREATE TABLE IF NOT EXISTS recordings (
    id bigserial PRIMARY KEY,
    file_path text NOT NULL,
    owner_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    title TEXT NOT NULL
);