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
	mux.Get("/auth/google", http.HandlerFunc(app.signUpWithGoogleProvider))
	mux.Get("/auth/google/callback", http.HandlerFunc(app.signUpWithGoogleCallback))

	return mux
}
