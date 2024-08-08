package data

import (
	"time"
)

type Tune struct {
	ID                 int64     `json:"id"`
	CreatedAt          time.Time `json:"created_at"`
	Version            int32     `json:"version"`
	Title              string    `json:"title"`
	Keys               []Key     `json:"keys"`
	TimeSignatureUpper int8      `json:"time_signature_upper"`
	TimeSignatureLower int8      `json:"time_signature_lower"`
	BandID             int64     `json:"band_id"`
	Status             string    `json:"status"`
}
