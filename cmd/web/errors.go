package main

import "net/http"

func (app *application) serverError(w http.ResponseWriter, err error) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func (app *application) NotFound(w http.ResponseWriter) {
	http.Error(w, "Not found", http.StatusNotFound)
}
