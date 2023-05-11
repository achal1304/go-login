package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) Routes() http.Handler {
	mux := pat.New()

	// mux.Get("/", http.HandlerFunc(Home))
	mux.Get("/", app.session.Enable(http.HandlerFunc(app.signUpUserForm)))
	mux.Post("/", app.session.Enable(http.HandlerFunc(app.signUpUser)))
	mux.Get("/login", app.session.Enable(http.HandlerFunc(app.loginUserForm)))
	mux.Post("/login", app.session.Enable(http.HandlerFunc(app.loginUser)))
	mux.Get("/auth/google", app.session.Enable(http.HandlerFunc(app.signUpWithGoogleProvider)))
	mux.Get("/auth/google/callback", app.session.Enable(http.HandlerFunc(app.signUpWithGoogleCallback)))
	mux.Get("/home/:id", app.session.Enable(http.HandlerFunc(app.profile)))
	mux.Get("/home/edit/:id", app.session.Enable(http.HandlerFunc(app.editProfileForm)))
	mux.Patch("/home/edit/:id", app.session.Enable(http.HandlerFunc(app.editProfile)))

	return app.methodOverride(mux)
}
