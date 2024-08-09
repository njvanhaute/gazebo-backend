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
	query := `
		INSERT INTO bands (name)
		VALUES ($1)
		RETURNING id, created_at, version`

	args := []any{band.Name}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return b.DB.QueryRowContext(ctx, query, args...).Scan(&band.ID, &band.CreatedAt, &band.Version)
}

func (b BandModel) Get(id int64) (*Band, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, version, name
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
		SET name = $1, version = version + 1
		WHERE id = $2 AND version = $3
		RETURNING version`

	args := []any{
		band.Name,
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
