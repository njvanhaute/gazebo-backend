package main

import (
	"errors"
	"fmt"
	"net/http"

	"gazebo.njvanhaute.com/internal/data"
	"gazebo.njvanhaute.com/internal/validator"
)

func (app *application) createBandHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	band := &data.Band{
		Name:    input.Name,
		OwnerID: user.ID,
	}

	v := validator.New()

	if data.ValidateBand(v, band); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Bands.Insert(band)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/bands/%d", band.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"band": band}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getBandHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	band, err := app.models.Bands.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"band": band}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getMyBandsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	bands, err := app.models.BandMembers.GetAllBandsForUser(user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"bands": bands}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateBandHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	band, err := app.models.Bands.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user := app.contextGetUser(r)

	if band.OwnerID != user.ID {
		app.notPermittedResponse(w, r)
		return
	}

	var input struct {
		Name *string `json:"name"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		band.Name = *input.Name
	}

	v := validator.New()

	if data.ValidateBand(v, band); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Bands.Update(band)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"band": band}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteBandHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	band, err := app.models.Bands.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user := app.contextGetUser(r)

	if band.OwnerID != user.ID {
		app.notPermittedResponse(w, r)
		return
	}

	err = app.models.Bands.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "band successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) addUserToBandHandler(w http.ResponseWriter, r *http.Request) {
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
	headers.Set("Location", fmt.Sprintf("/v1/bands/%d/users/%d", member.BandID, member.UserID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"member": member}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) removeUserFromBandHandler(w http.ResponseWriter, r *http.Request) {
	bandID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	userToBeDeletedID, err := app.readIntParam("userId", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	band, err := app.models.Bands.Get(bandID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	requestingUser := app.contextGetUser(r)

	if requestingUser.ID != userToBeDeletedID && requestingUser.ID != band.OwnerID {
		app.notPermittedResponse(w, r)
		return
	}

	if userToBeDeletedID == band.OwnerID {
		app.cannotRemoveOwnerResponse(w, r)
		return
	}

	member := data.BandMember{
		BandID: bandID,
		UserID: userToBeDeletedID,
	}

	err = app.models.BandMembers.Delete(&member)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "user successfully removed"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
