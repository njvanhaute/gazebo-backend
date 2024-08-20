package data

import (
	"database/sql"
	"time"
)

type Recording struct {
	ID        int64     `json:"id"`
	OwnerID   int64     `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	FilePath  string    `json:"-"`
	FileType  string    `json:"file_type"`
	Title     string    `json:"title"`
}

type RecordingModel struct {
	DB *sql.DB
}
