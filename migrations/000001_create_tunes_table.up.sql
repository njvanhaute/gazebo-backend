CREATE TYPE tunestatus as ENUM ('germinating', 'seedling', 'flowering');

CREATE TABLE IF NOT EXISTS tunes (
    id bigserial PRIMARY KEY,
    version integer NOT NULL DEFAULT 1,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    keys text[] NOT NULL,
    time_signature_upper integer NOT NULL,
    time_signature_lower integer NOT NULL,
    status tunestatus NOT NULL
);