package main

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAuthConfig = &oauth2.Config{
	ClientID:     "506670931401-ud0mqe2gf16o06r7llks4bg9a19tohck.apps.googleusercontent.com",
	ClientSecret: "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",
	RedirectURL:  "http://localhost:3000/auth/google/callback",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}
