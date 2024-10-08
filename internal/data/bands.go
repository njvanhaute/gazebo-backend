package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"gazebo.njvanhaute.com/internal/validator"
)

type Band struct {
	ID        int64     `json:"id"`
	OwnerID   int64     `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	Version   int32     `json:"version"`
	Name      string    `json:"name"`
}

func ValidateBand(v *validator.Validator, band *Band) {
	v.Check(band.Name != "", "name", "must be provided")
	v.Check(len(band.Name) <= 500, "name", "must not be more than 500 bytes long")
}

type BandModel struct {
	DB *sql.DB
}

func (b BandModel) Insert(band *Band) error {
	insertBandQuery := `
		INSERT INTO bands (name, owner_id)
		VALUES ($1, $2)
		RETURNING id, created_at, version`

	args := []any{band.Name, band.OwnerID}

	tx, err := b.DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = tx.QueryRowContext(ctx, insertBandQuery, args...).Scan(&band.ID, &band.CreatedAt, &band.Version)
	if err != nil {
		return err
	}

	addOwnerToBandQuery := `
		INSERT INTO band_members (band_id, user_id)
		VALUES ($1, $2)`

	args = []any{band.ID, band.OwnerID}

	_, err = tx.ExecContext(ctx, addOwnerToBandQuery, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (b BandModel) Get(id int64) (*Band, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, version, name, owner_id
		FROM bands
		WHERE id = $1`

	var band Band

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := b.DB.QueryRowContext(ctx, query, id).Scan(
		&band.ID,
		&band.CreatedAt,
		&band.Version,
		&band.Name,
		&band.OwnerID,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &band, nil
}

func (b BandModel) Update(band *Band) error {
	query := `
		UPDATE bands
		SET name = $1, owner_id = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version`

	args := []any{
		band.Name,
		band.OwnerID,
		band.ID,
		band.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := b.DB.QueryRowContext(ctx, query, args...).Scan(&band.Version)
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

func (b BandModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM bands
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := b.DB.ExecContext(ctx, query, id)
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
