package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) readIntParam(name string, r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	intParam, err := strconv.ParseInt(params.ByName(name), 10, 64)
	if err != nil || intParam < 1 {
		return 0, fmt.Errorf("invalid int passed to parameter \"%s\"", name)
	}

	return intParam, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
