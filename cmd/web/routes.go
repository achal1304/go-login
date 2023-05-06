package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) Routes() http.Handler {
	mux := pat.New()

	// mux.Get("/", http.HandlerFunc(Home))
	mux.Get("/", http.HandlerFunc(app.signUpUserForm))
	mux.Post("/", http.HandlerFunc(app.signUpUser))

	return mux
}
