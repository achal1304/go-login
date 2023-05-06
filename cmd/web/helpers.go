package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, files []string, data interface{}) {
	ts, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Print(err)
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		fmt.Print(err)
		app.serverError(w, err)
	}
}
