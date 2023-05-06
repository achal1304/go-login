package main

import (
	"github.com/achal1304/go-login/pkg/forms"
	"github.com/achal1304/go-login/pkg/models"
)

type templateData struct {
	Users *models.User
	Form  *forms.Form
}
