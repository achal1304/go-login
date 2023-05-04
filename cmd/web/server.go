package main

import (
	"log"
	"net/http"
	"time"
)

func RunServer() {
	srv := &http.Server{
		Addr:         ":4000",
		Handler:      Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := srv.ListenAndServe()
	log.Fatal(err)
}
