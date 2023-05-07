package main

import (
	"log"
	"net/http"
	"time"
)

func (app *application) RunServer() {
	srv := &http.Server{
		Addr:         ":3000",
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := srv.ListenAndServe()
	log.Fatal(err)
}
