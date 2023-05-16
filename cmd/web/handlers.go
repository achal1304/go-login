package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/achal1304/go-login/internal/data"
	"github.com/achal1304/go-login/pkg/forms"
	"github.com/achal1304/go-login/pkg/models"
	"github.com/achal1304/go-login/pkg/models/mysql"
	"golang.org/x/crypto/bcrypt"
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
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/login.page.tmpl",
	}
	data := templateData{
		Form: forms.New(nil),
	}
	app.render(w, r, files, data)
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.MinLength("password", 3)
	form.MatchesPattern("email", forms.EmailRX)

	if !form.Valid() {
		app.render(w, r, []string{"./ui/html/login.page.tmpl"}, &templateData{Form: form})
		return
	}

	userId, err := app.users.AuthenticateUser(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, mysql.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, []string{"./ui/html/login.page.tmpl"}, &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "authenticatedUserID", userId)

	http.Redirect(w, r, fmt.Sprintf("/home/%d", userId), http.StatusSeeOther)
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

	userId, err := app.users.GetUser(user.Email)
	if err != nil {
		app.serverError(w, err)
	}
	if userId != 1 {
		app.session.Put(r, "authenticatedUserID", userId)
		http.Redirect(w, r, fmt.Sprintf("/home/%d", userId), http.StatusSeeOther)
	} else {
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
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}
}

func (app *application) profile(w http.ResponseWriter, r *http.Request) {
	// userId := app.session.Get(r, "authenticatedUserID")
	// w.Write([]byte(fmt.Sprintf("Welcome to home page %s", userId)))

	userId, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || userId < 0 {
		fmt.Println(err)
		app.NotFound(w)
		return
	}

	user, err := app.users.GetUserDetailsFromId(userId)
	if err != nil {
		fmt.Println(err)
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
	}
	app.render(w, r, files, &templateData{User: user})

}

func (app *application) editProfileForm(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || userId < 0 {
		fmt.Println(err)
		app.NotFound(w)
		return
	}
	fmt.Print("ID - edit ", userId)

	user, err := app.users.GetUserDetailsFromId(userId)
	if err != nil {
		fmt.Println(err)
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html/editprofile.page.tmpl",
	}

	app.render(w, r, files, &templateData{Form: forms.New(nil), User: user})
}

func (app *application) editProfile(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.Atoi(r.URL.Query().Get(":id"))
	fmt.Print("ID - patch -  ", userId)

	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("name", 3)
	form.MinLength("address", 5)

	user, err := app.users.GetUserDetailsFromId(userId)
	if err != nil {
		fmt.Println(err)
		app.serverError(w, err)
		return
	}
	err = bcrypt.CompareHashAndPassword(user.Hashed_password, []byte(defaultPassword))
	if err == nil && user.Email != form.Get("email") {
		form.Errors.Add("email", "This field cannot be changed as you are signedIn from Google")
	}

	if !form.Valid() {
		fmt.Println("Invalid form")
		app.render(w, r, []string{"./ui/html/editprofile.page.tmpl"}, &templateData{Form: form, User: user})
		return
	}

	err = app.users.UpdateUserDetails(userId, form.Get("email"), form.Get("address"), form.Get("name"))
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, mysql.ErrDuplicateEmail) {
			form.Errors.Add("email", "Email already exists")
			app.render(w, r, []string{"./ui/html/editprofile.page.tmpl"}, &templateData{Form: form, User: user})
		} else {
			app.serverError(w, err)
		}
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/home/%d", userId), http.StatusSeeOther)
}

func (app *application) resetPasswordForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside resetPasswordForm")
	files := []string{
		"./ui/html/resetpassword.page.tmpl",
	}
	data := templateData{
		Form: forms.New(nil),
	}
	app.render(w, r, files, data)
}

func (app *application) resetPassword(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.MatchesPattern("email", forms.EmailRX)
	if !form.Valid() {
		fmt.Println("Invalid form")
		app.render(w, r, []string{"./ui/html/resetpassword.page.tmpl"}, &templateData{Form: form})
		return
	}

	isEmailPresent := app.users.IsEmailPresent(form.Get("email"))
	if !isEmailPresent {
		form.Errors.Add("email", "Email doesnot exist, please enter a valid email")
		app.render(w, r, []string{"./ui/html/resetpassword.page.tmpl"}, &templateData{Form: form})
		return
	}

	user, err := app.users.GetUser(form.Get("email"))
	fmt.Println("user id is ", user)
	if err != nil {
		app.serverError(w, err)
	}

	err = app.mailer.Send(form.Get("email"), "resetlink.tmpl", user, form.Get("email"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) resetPasswordWithToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get(":token")

	claims, err := data.DecodeJWT(token)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Invalid token"))
		return
	}

	fmt.Println(claims.UserId, claims.Email)
}
