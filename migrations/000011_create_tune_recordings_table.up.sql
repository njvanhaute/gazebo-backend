CREATE TABLE IF NOT EXISTS tune_recordings (
    tune_id bigint REFERENCES tunes ON DELETE CASCADE,
    recording_id bigint REFERENCES recordings ON DELETE CASCADE,
    CONSTRAINT tune_recording_pkey PRIMARY KEY (tune_id, recording_id)
);