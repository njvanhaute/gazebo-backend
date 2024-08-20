CREATE TABLE IF NOT EXISTS documents (
    id bigserial PRIMARY KEY,
    tune_id bigint NOT NULL REFERENCES tunes ON DELETE CASCADE,
    file_path text NOT NULL,
    owner_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    title TEXT NOT NULL
);