package data

import (
	"context"
	"database/sql"
	"time"

	"gazebo.njvanhaute.com/internal/validator"
)

type BandMember struct {
	BandID int64 `json:"band_id"`
	UserID int64 `json:"user_id"`
}
type BandMemberModel struct {
	DB *sql.DB
}

func ValidateBandMember(v *validator.Validator, member *BandMember) {
	v.Check(member.BandID > 0, "band_id", "must be a positive integer")
	v.Check(member.UserID > 0, "user_id", "must be a positive integer")
}

func (b BandMemberModel) Insert(member *BandMember) error {
	query := `
		INSERT INTO band_members (band_id, user_id)
		VALUES ($1, $2)`

	args := []any{member.BandID, member.UserID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := b.DB.ExecContext(ctx, query, args...)
	if err != nil {
		println(err.Error())

		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "band_user_pkey"`:
			return ErrRecordAlreadyExists
		case err.Error() == `pq: insert or update on table "band_members" violates foreign key constraint "band_members_user_id_fkey"`:
			return ErrUserNotFound
		case err.Error() == `pq: insert or update on table "band_members" violates foreign key constraint "band_members_band_id_fkey"`:
			return ErrBandNotFound
		default:
			return err
		}
	}

	return nil
}

func (b BandMemberModel) GetAllBandsForUser(id int64) ([]*Band, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT bands.id, bands.created_at, bands.version, bands.name
		FROM bands, band_members
		WHERE bands.id = band_members.band_id
		AND band_members.user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := b.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	bands := []*Band{}

	for rows.Next() {
		var band Band

		err := rows.Scan(
			&band.ID,
			&band.CreatedAt,
			&band.Version,
			&band.Name,
		)

		if err != nil {
			return nil, err
		}

		bands = append(bands, &band)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bands, nil
}

func (b BandMemberModel) UserIsInBand(userId int64, bandId int64) (bool, error) {
	if userId < 1 || bandId < 1 {
		return false, ErrRecordNotFound
	}

	query := `
		SELECT COUNT(*)
		FROM band_members
		WHERE user_id = $1 AND band_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	count := 0
	args := []any{userId, bandId}

	err := b.DB.QueryRowContext(ctx, query, args...).Scan(&count)

	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func (b BandMemberModel) GetAllUsersForBand(id int64) ([]*User, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT users.id, users.created_at, users.name, users.email, users.activated
		FROM users, band_members
		WHERE users.id = band_members.user_id
		AND band_members.band_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := b.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.CreatedAt,
			&user.Email,
			&user.Name,
			&user.Activated,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (b BandMemberModel) Delete(member *BandMember) error {
	if member.BandID < 1 || member.UserID < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM band_members
		WHERE band_id = $1
		AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := b.DB.ExecContext(ctx, query, member.BandID, member.UserID)
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
