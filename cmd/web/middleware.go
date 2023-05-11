package main

import (
	"fmt"
	"net/http"
)

func (app *application) methodOverride(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("in method override middleware")
		if r.Method == http.MethodPost {
			fmt.Println("changing method")
			if r.FormValue("_method") == "PATCH" {
				r.Method = http.MethodPatch
			}
		}
		next.ServeHTTP(w, r)
	})
}
