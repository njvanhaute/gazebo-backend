package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
		INSERT INTO tunes (title, keys, time_signature_upper, time_signature_lower, status, band_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, version`

	args := []any{tune.Title, pq.Array(tune.Keys), tune.TimeSignatureUpper, tune.TimeSignatureLower, tune.Status, tune.BandID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return t.DB.QueryRowContext(ctx, query, args...).Scan(&tune.ID, &tune.CreatedAt, &tune.Version)
}

func (t TuneModel) Get(id int64) (*Tune, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, version, title, keys, time_signature_upper, time_signature_lower, status, band_id
		FROM tunes
		WHERE id = $1`

	var tune Tune
	var keys []string

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, id).Scan(
		&tune.ID,
		&tune.CreatedAt,
		&tune.Version,
		&tune.Title,
		pq.Array(&keys),
		&tune.TimeSignatureUpper,
		&tune.TimeSignatureLower,
		&tune.Status,
		&tune.BandID,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	for _, key := range keys {
		tune.Keys = append(tune.Keys, Key(key))
	}

	return &tune, nil
}

func (t TuneModel) GetAll(bandId int64, title string, keys []string, statuses []string, filters Filters) ([]*Tune, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, version, title, keys, time_signature_upper, time_signature_lower, status, band_id
		FROM tunes
		WHERE band_id = $1
		AND (to_tsvector('simple', title) @@ plainto_tsquery('simple', $2) OR $2 = '')
		AND (keys @> $3 OR $3 = '{}')
		AND (status = ANY($4) or $4 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $5 OFFSET $6`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{bandId, title, pq.Array(keys), pq.Array(statuses), filters.limit(), filters.offset()}

	rows, err := t.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	tunes := []*Tune{}

	for rows.Next() {
		var tune Tune
		var keys []string

		err := rows.Scan(
			&totalRecords,
			&tune.ID,
			&tune.CreatedAt,
			&tune.Version,
			&tune.Title,
			pq.Array(&keys),
			&tune.TimeSignatureUpper,
			&tune.TimeSignatureLower,
			&tune.Status,
			&tune.BandID,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		for _, key := range keys {
			tune.Keys = append(tune.Keys, Key(key))
		}

		tunes = append(tunes, &tune)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return tunes, metadata, nil
}

func (t TuneModel) Update(tune *Tune) error {
	query := `
		UPDATE tunes
		SET title = $1, keys = $2, time_signature_upper = $3, time_signature_lower = $4, status = $5, version = version + 1
		WHERE id = $6 AND version = $7
		RETURNING version`

	args := []any{
		tune.Title,
		pq.Array(tune.Keys),
		tune.TimeSignatureUpper,
		tune.TimeSignatureLower,
		tune.Status,
		tune.ID,
		tune.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, args...).Scan(&tune.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (t TuneModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM tunes
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := t.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
