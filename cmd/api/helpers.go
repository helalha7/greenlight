package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type envelope map[string]any

func (app *application) readIdParam(r *http.Request) (int, error) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, statusCode int, data any, header http.Header) error {

	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, values := range header {
		for _, value := range values {
			w.Header().Set(key, value)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	w.Write(js)

	return nil
}
