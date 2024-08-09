package main

import (
	"errors"
	"fmt"
	"net/http"

	"gazebo.njvanhaute.com/internal/data"
	"gazebo.njvanhaute.com/internal/validator"
)

func (app *application) addBandMember(w http.ResponseWriter, r *http.Request) {
	bandID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		UserID int64 `json:"user_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	member := &data.BandMember{
		BandID: bandID,
		UserID: input.UserID,
	}

	v := validator.New()

	if data.ValidateBandMember(v, member); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.BandMembers.Insert(member)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordAlreadyExists):
			v.AddError("user_id", "this user is already in this band")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrUserNotFound):
			v.AddError("user_id", "this user does not exist")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrBandNotFound):
			v.AddError("band_id", "this band does not exist")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/bands/%d/members/%d", member.BandID, member.UserID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"band_member": member}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) getBandMembers(w http.ResponseWriter, r *http.Request) {
	bandID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Make sure we have a real band ID
	_, err = app.models.Bands.Get(bandID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	members, err := app.models.BandMembers.GetAllUsersForBand(bandID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"band_id": bandID, "members": members}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteBandMember(w http.ResponseWriter, r *http.Request) {
	// TODO
	bandId, err := app.readIDParam(r)
	if err != nil {
		panic("helP")
	}

	userId, err := app.readIntParam("userId", r)
	if err != nil {
		panic("help")
	}

	print(bandId, userId)
}
