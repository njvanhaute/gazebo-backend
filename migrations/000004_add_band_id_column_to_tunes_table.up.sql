ALTER TABLE tunes ADD COLUMN band_id int NOT NULL;
ALTER TABLE tunes ADD CONSTRAINT fk_tune_band FOREIGN KEY (band_id) REFERENCES bands(id)