package main

import (
	"net/http"

	"gazebo.njvanhaute.com/internal/data"
	"gazebo.njvanhaute.com/internal/validator"
)

func (app *application) addBandMember(w http.ResponseWriter, r *http.Request) {
	var input struct {
		BandID int64 `json:"band_id"`
		UserID int64 `json:"user_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	member := &data.BandMember{
		BandID: input.BandID,
		UserID: input.UserID,
	}

	v := validator.New()

	if data.ValidateBandMember(v, member); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.BandMembers.Insert(member)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)

	err = app.writeJSON(w, http.StatusCreated, envelope{"band_member": member}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
