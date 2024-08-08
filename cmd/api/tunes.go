package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"gazebo.njvanhaute.com/internal/data"
	"gazebo.njvanhaute.com/internal/validator"
)

func (app *application) createTuneHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title              string     `json:"title"`
		Keys               []data.Key `json:"keys"`
		TimeSignatureUpper int8       `json:"time_signature_upper"`
		TimeSignatureLower int8       `json:"time_signature_lower"`
		BandID             int64      `json:"band_id"`
		Status             string     `json:"status"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tune := &data.Tune{
		Title:              input.Title,
		Keys:               input.Keys,
		TimeSignatureUpper: input.TimeSignatureUpper,
		TimeSignatureLower: input.TimeSignatureLower,
		BandID:             input.BandID,
		Status:             input.Status,
	}

	v := validator.New()

	if data.ValidateTune(v, tune); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Tunes.Insert(tune)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/tunes/%d", tune.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"tune": tune}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getTuneHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	tune, err := app.models.Tunes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tune": tune}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateTuneHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	tune, err := app.models.Tunes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title              *string    `json:"title"`
		Keys               []data.Key `json:"keys"`
		TimeSignatureUpper *int8      `json:"time_signature_upper"`
		TimeSignatureLower *int8      `json:"time_signature_lower"`
		Status             *string    `json:"status"`
		BandID             *int64     `json:"band_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		tune.Title = *input.Title
	}

	if input.Keys != nil {
		tune.Keys = input.Keys
	}

	if input.TimeSignatureUpper != nil {
		tune.TimeSignatureUpper = *input.TimeSignatureUpper
	}

	if input.TimeSignatureLower != nil {
		tune.TimeSignatureLower = *input.TimeSignatureLower
	}

	if input.Status != nil {
		tune.Status = *input.Status
	}

	if input.BandID != nil {
		tune.BandID = *input.BandID
	}

	// TODO: REMOVE THIS
	tune.BandID = 1

	v := validator.New()

	if data.ValidateTune(v, tune); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Tunes.Update(tune)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tune": tune}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getTunesForBandHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	tune := data.Tune{
		ID:                 id,
		CreatedAt:          time.Now(),
		Version:            1,
		Title:              "Pickin' Sage",
		Keys:               []data.Key{"F major"},
		TimeSignatureUpper: 4,
		TimeSignatureLower: 4,
		BandID:             1,
		Status:             "germinating",
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tune": tune}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTuneHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Tunes.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "tune successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
