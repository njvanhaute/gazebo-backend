package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"gazebo.njvanhaute.com/internal/validator"
)

type Document struct {
	ID        int64     `json:"id"`
	TuneID    int64     `json:"tune_id"`
	OwnerID   int64     `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	FilePath  string    `json:"-"`
	FileType  string    `json:"file_type"`
	Title     string    `json:"title"`
}

type DocumentModel struct {
	DB *sql.DB
}

func ValidateDocument(v *validator.Validator, doc *Document) {
	v.Check(doc.FileType != "", "file_type", "must be provided")
	validFileTypes := []string{"pdf"}
	v.Check(validator.PermittedValue(doc.FileType, validFileTypes...), "file_type", "invalid file type")

	v.Check(doc.Title != "", "title", "must be provided")
	v.Check(len(doc.Title) <= 500, "title", "must not be more than 500 bytes long")
}

func (d DocumentModel) Get(id int64) (*Document, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, tune_id, owner_id, created_at, file_path, file_type, title
		FROM documents
		WHERE id = $1`

	var doc Document

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := d.DB.QueryRowContext(ctx, query, id).Scan(
		&doc.ID,
		&doc.TuneID,
		&doc.OwnerID,
		&doc.CreatedAt,
		&doc.FilePath,
		&doc.FileType,
		&doc.Title,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &doc, nil
}

func (d DocumentModel) Insert(doc *Document) error {
	query := `
		INSERT INTO documents (tune_id, owner_id, file_path, file_type, title)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	args := []any{doc.TuneID, doc.OwnerID, doc.FilePath, doc.FileType, doc.Title}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return d.DB.QueryRowContext(ctx, query, args...).Scan(&doc.ID, &doc.CreatedAt)
}

func (d DocumentModel) GetAllDocsForTune(tuneID int64) ([]*Document, error) {
	query := `
		SELECT id, tune_id, owner_id, created_at, file_path, file_type, title
		FROM documents
		WHERE tune_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := d.DB.QueryContext(ctx, query, tuneID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	docs := []*Document{}

	for rows.Next() {
		var doc Document

		err := rows.Scan(
			&doc.ID,
			&doc.TuneID,
			&doc.OwnerID,
			&doc.CreatedAt,
			&doc.FilePath,
			&doc.FileType,
			&doc.Title,
		)

		if err != nil {
			return nil, err
		}

		docs = append(docs, &doc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return docs, nil
}
