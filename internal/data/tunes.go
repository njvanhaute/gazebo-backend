package data

import (
	"database/sql"
	"time"

	"gazebo.njvanhaute.com/internal/validator"
	"github.com/lib/pq"
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

func ValidateTune(v *validator.Validator, tune *Tune) {
	v.Check(tune.Title != "", "title", "must be provided")
	v.Check(len(tune.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(tune.Keys != nil, "keys", "must be provided")
	v.Check(len(tune.Keys) >= 1, "keys", "must contain at least 1 key")
	v.Check(validator.Unique(tune.Keys), "keys", "must not contain duplicate values")

	v.Check(tune.TimeSignatureUpper != 0, "time_signature_upper", "must be provided")
	v.Check(tune.TimeSignatureUpper > 1, "time_signature_upper", "must be a positive integer")

	v.Check(tune.TimeSignatureLower != 0, "time_signature_lower", "must be provided")
	v.Check(tune.TimeSignatureLower > 1, "time_signature_lower", "must be at least 2")
	v.Check(tune.TimeSignatureLower&(tune.TimeSignatureLower-1) == 0, "time_signature_lower", "must be a power of 2")

	v.Check(tune.BandID != 0, "band_id", "must be provided")
	v.Check(tune.BandID > 0, "band_id", "must be a positive integer")

	validStatusList := []string{"germinating", "seedling", "flowering"}
	v.Check(validator.PermittedValue(tune.Status, validStatusList...), "status", "invalid status value")

}

type TuneModel struct {
	DB *sql.DB
}

func (t TuneModel) Insert(tune *Tune) error {
	query := `
		INSERT INTO tunes (title, keys, time_signature_upper, time_signature_lower, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, version`

	args := []any{tune.Title, pq.Array(tune.Keys), tune.TimeSignatureUpper, tune.TimeSignatureLower, tune.Status}
	return t.DB.QueryRow(query, args...).Scan(&tune.ID, &tune.CreatedAt, &tune.Version)
}

func (t TuneModel) GetAllForBand(bandId int64) ([]*Tune, error) {
	return nil, nil
}

func (t TuneModel) CreateNewVersion(tune *Tune) error {
	return nil
}

func (t TuneModel) DeleteVersion(id int64, version int32) error {
	return nil
}

func (t TuneModel) DeleteTune(id int64) error {
	return nil
}
