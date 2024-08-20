package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"gazebo.njvanhaute.com/internal/data"
	"github.com/google/uuid"
)

func (app *application) documentUploadHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TuneID   int64  `json:"tune_id"`
		OwnerID  int64  `json:"owner_id"`
		FileType string `json:"file_type"`
		Title    string `json:"title"`
	}

	mr, err := r.MultipartReader()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	id := uuid.New()
	filePath := "./docs/" + id.String()

	for {
		part, err := mr.NextPart()

		if err == io.EOF {
			break
		}

		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		if part.FormName() == "info" {
			dec := json.NewDecoder(part)
			dec.DisallowUnknownFields()

			err = dec.Decode(&input)
			if err != nil {
				app.badRequestResponse(w, r, handleJSONDecodingErrors(err))
				return
			}
		}

		if part.FormName() == "file" {
			outfile, err := os.Create(filePath)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			defer outfile.Close()

			_, err = io.Copy(outfile, part)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
		}
	}

	doc := &data.Document{
		TuneID:   input.TuneID,
		OwnerID:  input.OwnerID,
		FilePath: filePath,
		FileType: input.FileType,
		Title:    input.Title,
	}

	// TODO: validation, insert into DB, return response

}
