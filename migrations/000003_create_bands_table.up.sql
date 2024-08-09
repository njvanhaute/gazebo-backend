CREATE TABLE IF NOT EXISTS bands (
    id bigserial PRIMARY KEY,
    version integer NOT NULL DEFAULT 1,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL
);