package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// Tunes
	router.HandlerFunc(http.MethodGet, "/v1/tunes/:id", app.requireActivatedUser(app.getTuneHandler))
	router.HandlerFunc(http.MethodPost, "/v1/tunes", app.requireActivatedUser(app.createTuneHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/tunes/:id", app.requireActivatedUser(app.updateTuneHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/tunes/:id", app.requireActivatedUser(app.deleteTuneHandler))

	// Bands
	router.HandlerFunc(http.MethodGet, "/v1/bands/:id", app.requireActivatedUser(app.getBandHandler))
	router.HandlerFunc(http.MethodPost, "/v1/bands", app.requireActivatedUser(app.createBandHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/bands/:id", app.requireActivatedUser(app.updateBandHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/bands/:id", app.requireActivatedUser(app.deleteBandHandler))
	router.HandlerFunc(http.MethodGet, "/v1/bands/:id/tunes/", app.requireActivatedUser(app.listTunesForBandHandler))

	// Band members
	router.HandlerFunc(http.MethodPost, "/v1/bands/:id/users", app.requireActivatedUser(app.addUserToBandHandler))
	router.HandlerFunc(http.MethodGet, "/v1/bands/:id/users", app.requireActivatedUser(app.getUsersInBandHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/bands/:id/users/:userId", app.requireActivatedUser(app.removeUserFromBandHandler))
	router.HandlerFunc(http.MethodGet, "/v1/users/:id/bands", app.requireActivatedUser(app.getBandsJoinedByUserHandler))

	// Users
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	// Authentication
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
