package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"

	"log"
)

func main() {

	//Client secret created on google cloud platform/ Apis & Services / Credentials
	key := "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"

	//Time period over which the token is valid(or existant)
	maxAge := 86400
	isProd := false

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"

	//On deafult should be enabled
	store.Options.HttpOnly = true
	//Enables or disables https protocol
	store.Options.Secure = isProd

	gothic.Store = store

	//Creates provider for google using Client id and Client secret
	goth.UseProviders(
		google.New("261473284823-sh61p2obchbmdrq9pucc7s5oo9c8l98j.apps.googleusercontent.com", "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", "http://localhost:3000/auth/google/callback", "email", "profile"),
	)

	//Http router for Go
	p := pat.New()
	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}

		// In success.html file or any front display we can select what we want to display.
		// user is instance of goth.User class of which fields can be found in documentation.
		t, _ := template.ParseFiles("pages/success.html")
		t.Execute(res, user)
	})
	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.ParseFiles("pages/index.html")
		t.Execute(res, false)
	})

	log.Println("listening on localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", p))
}
