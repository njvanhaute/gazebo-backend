ALTER TABLE tunes ADD CONSTRAINT tunes_time_signature_upper_check CHECK(time_signature_upper >= 1);
ALTER TABLE tunes ADD CONSTRAINT tunes_time_signature_lower_check CHECK(time_signature_lower >= 2);