CREATE TABLE IF NOT EXISTS recordings (
    id bigserial PRIMARY KEY,
    file_path text NOT NULL,
    file_type text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    owner_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    title TEXT NOT NULL
);