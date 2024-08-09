CREATE TABLE IF NOT EXISTS band_members (
    band_id int REFERENCES bands ON DELETE CASCADE,
    user_id int REFERENCES users ON DELETE CASCADE,
    CONSTRAINT band_user_pkey PRIMARY KEY (band_id, user_id)
);