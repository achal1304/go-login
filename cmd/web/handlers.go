package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/achal1304/go-login/pkg/forms"
	"github.com/achal1304/go-login/pkg/models"
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

func (app *application) signUpWithGoogleProvider(w http.ResponseWriter, r *http.Request) {
	url := googleOAuthConfig.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *application) signUpWithGoogleCallback(w http.ResponseWriter, r *http.Request) {
	token, err := googleOAuthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		fmt.Errorf("failed to exchange code for token: %s", err.Error())
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		fmt.Errorf("Get: " + err.Error() + "\n")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	var user models.GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		fmt.Printf("Error in decoding")
	}

	w.Write([]byte("Hello, I'm protected\n"))
	w.Write([]byte(string(user.Email)))

	form := forms.New(r.PostForm)
	err = app.users.Insert(user.Email, defaultPassword)
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

	return
}
