package main

import (
	"fmt"
	"net/http"

	"github.com/achal1304/go-login/pkg/forms"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello home"))
	// fmt.Print("Hiiiii")
}

func (app *application) signUpUserForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside signupuserform")
	files := []string{
		"./ui/html/signup.page.tmpl",
	}
	data := templateData{
		Form: forms.New(nil),
	}
	app.render(w, r, files, data)
}
func (app *application) signUpUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is signup user form"))
}
