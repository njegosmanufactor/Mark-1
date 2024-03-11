package LoginRegister

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"

	applicationService "MongoGoogle/ApplicationService"
	gitService "MongoGoogle/GitService"
	googleService "MongoGoogle/GoogleService"
	userType "MongoGoogle/Model"
)

func Authentication() {
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

	mux := mux.NewRouter()

	//Homepage display on path "https://localhost:3000"
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.ParseFiles("LoginRegister/pages/index.html")
		t.Execute(res, false)
	})

	//Using google OAuth2 to authenticate user
	mux.HandleFunc("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	mux.HandleFunc("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
		googleService.CompleteGoogleUserAuthentication(res, req)
	})

	//Github authentication paths
	mux.HandleFunc("/login/github", func(w http.ResponseWriter, r *http.Request) {
		gitService.GithubLoginHandler(w, r)
	})

	// Github callback
	mux.HandleFunc("/login/github/callback", func(w http.ResponseWriter, r *http.Request) {
		gitService.GithubCallbackHandler(w, r)
	})

	// Route where the authenticated user is redirected to
	mux.HandleFunc("/loggedin", func(w http.ResponseWriter, r *http.Request) {
		var nill userType.GitHubData
		gitService.LoggedinHandler(w, r, nill)
	})

	//Register page display
	mux.HandleFunc("/register.html", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("LoginRegister/pages/register.html")
		if err != nil {
			fmt.Fprintf(res, "Error parsing template: %v", err)
			return
		}
		t.Execute(res, false)
	})

	//Our service that serves login functionality
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		applicationService.ApplicationLogin(w, r)
	})

	//Our service that serves registration functionality
	mux.HandleFunc("/register", func(res http.ResponseWriter, req *http.Request) {

		applicationService.ApplicationRegister(res, req)

	})

	//Mux router listens for requests on port : 3000
	fmt.Println("[ UP ON PORT 3000 ]")
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
