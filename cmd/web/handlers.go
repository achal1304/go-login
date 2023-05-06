package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/achal1304/go-login/pkg/forms"
	"github.com/achal1304/go-login/pkg/models/mysql"
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
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.MinLength("password", 3)
	form.MatchesPattern("email", forms.EmailRX)

	if !form.Valid() {
		app.render(w, r, []string{"./ui/html/signup.page.tmpl"}, &templateData{Form: form})
		return
	}

	err = app.users.Insert(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, mysql.ErrDuplicateEmail) {
			form.Errors.Add("email", "Email already exists")
			app.render(w, r, []string{"./ui/html/signup.page.tmpl"}, &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}
	fmt.Print("Created a new user")
}
