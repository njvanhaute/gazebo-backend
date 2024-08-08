package main

import (
	"fmt"
	"net/http"
	"time"

	"gazebo.njvanhaute.com/internal/data"
)

func (app *application) createTuneHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new tune")
}

func (app *application) getTunesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIntParam("bandId", r)
	if err != nil {
		print(err.Error())
		http.NotFound(w, r)
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

	err = app.writeJSON(w, http.StatusOK, tune, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
