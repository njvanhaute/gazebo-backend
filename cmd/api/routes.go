package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tunes", app.createTuneHandler)
	router.HandlerFunc(http.MethodGet, "/v1/tunes/:bandId", app.getTunesHandler)

	return router
}
