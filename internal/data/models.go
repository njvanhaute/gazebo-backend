package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordAlreadyExists = errors.New("record already exists")
	ErrRecordNotFound      = errors.New("record not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrBandNotFound        = errors.New("band not found")
	ErrEditConflict        = errors.New("edit conflict")
)

type Models struct {
	BandMembers BandMemberModel
	Bands       BandModel
	Tokens      TokenModel
	Tunes       TuneModel
	Users       UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		BandMembers: BandMemberModel{DB: db},
		Bands:       BandModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Tunes:       TuneModel{DB: db},
		Users:       UserModel{DB: db},
	}
}
