package main

import (
	"fmt"
	"net/http"
	"time"

	"gazebo.njvanhaute.com/internal/data"
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

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) getTunesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIntParam("bandId", r)
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
