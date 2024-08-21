package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"gazebo.njvanhaute.com/internal/data"
	"gazebo.njvanhaute.com/internal/validator"
	"github.com/google/uuid"
)

func (app *application) documentUploadHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TuneID   int64  `json:"tune_id"`
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
	var doc *data.Document

	gotFile, gotMetadata := false, false
	numParts := 0

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

			user := app.contextGetUser(r)

			doc = &data.Document{
				TuneID:   input.TuneID,
				OwnerID:  user.ID,
				FilePath: filePath,
				FileType: input.FileType,
				Title:    input.Title,
			}

			v := validator.New()

			if data.ValidateDocument(v, doc); !v.Valid() {
				app.failedValidationResponse(w, r, v.Errors)
				return
			}

			tune, err := app.models.Tunes.Get(doc.TuneID)
			if err != nil {
				switch {
				case errors.Is(err, data.ErrRecordNotFound):
					v.AddError("tune_id", "invalid tune ID supplied")
					app.failedValidationResponse(w, r, v.Errors)
				default:
					app.serverErrorResponse(w, r, err)
				}
				return
			}

			userIsInBand, err := app.models.BandMembers.UserIsInBand(user.ID, tune.BandID)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			if !userIsInBand {
				v.AddError("tune_id", "you are not in the band that owns this tune")
				app.failedValidationResponse(w, r, v.Errors)
				return
			}

			gotMetadata = true
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

			gotFile = true
		}

		numParts += 1
	}

	if numParts != 2 {
		app.wrongNumberOfPartsResponse(w, r)
		return
	}

	if !gotFile {
		app.missingFileResponse(w, r)
		return
	}

	if !gotMetadata {
		app.missingMetadataResponse(w, r)
		return
	}

	err = app.models.Documents.Insert(doc)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		err = os.Remove(filePath)
		if err != nil {
			app.logger.Error(err.Error())
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/documents/%d", doc.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"doc": doc}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
