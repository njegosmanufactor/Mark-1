package LoginRegister

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"log"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"

	data "MongoGoogle/MongoDB"
)

// Function called before main initialisation that loads env variables for github oauth

/*func Authentication() {

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
	p := mux.NewRouter()

	p.HandleFunc("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}
		t, _ := template.ParseFiles("LoginRegister/pages/success.html")
		t.Execute(res, user)
		data.SaveUserGoogle(user.Email, user.FirstName, user.LastName)
	})
	p.HandleFunc("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	p.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("LoginRegister/pages/index.html")
		if err != nil {
			fmt.Fprintf(res, "Error parsing template: %v", err)
			return
		}
		t.Execute(res, false)
	})

	p.HandleFunc("/register.html", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("LoginRegister/pages/register.html")
		if err != nil {
			fmt.Fprintf(res, "Error parsing template: %v", err)
			return
		}
		t.Execute(res, false)
	})

	p.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
			return
		}

		// Čitanje podataka iz forme
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Poziv funkcije za proveru korisnika

		var valid bool = data.ValidUser(username, password)
		if valid {
			fmt.Fprintf(w, "Da")
		} else {
			fmt.Fprintf(w, "Ne")
		}

	})

	// Obrada podataka iz forme kada je ruta /register
	p.HandleFunc("/register", func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodPost {
			http.Error(res, "Only POST method allowed", http.StatusMethodNotAllowed)
			return
		}

		// Čitanje podataka iz forme
		username := req.FormValue("username")
		password := req.FormValue("password")

		// Poziv funkcije za čuvanje korisnika
		data.SaveUserApplication(username, password)

		// Redirekcija na neku drugu stranicu nakon registracije
	})

	log.Println("listening on localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", p))
}*/

// ////////////////////////////////////////////////////////////////////////////////////////////////
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
	mux.HandleFunc("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}
		t, _ := template.ParseFiles("LoginRegister/pages/success.html")
		t.Execute(res, user)
	})

	mux.HandleFunc("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.ParseFiles("LoginRegister/pages/index.html")
		t.Execute(res, false)
	})

	// Login route
	mux.HandleFunc("/login/github", func(w http.ResponseWriter, r *http.Request) {
		githubLoginHandler(w, r)
	})

	// Github callback
	mux.HandleFunc("/login/github/callback", func(w http.ResponseWriter, r *http.Request) {
		githubCallbackHandler(w, r)
	})

	// Route where the authenticated user is redirected to
	mux.HandleFunc("/loggedin", func(w http.ResponseWriter, r *http.Request) {
		loggedinHandler(w, r, "")
	})

	mux.HandleFunc("/register.html", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("LoginRegister/pages/register.html")
		if err != nil {
			fmt.Fprintf(res, "Error parsing template: %v", err)
			return
		}
		t.Execute(res, false)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
			return
		}

		// Čitanje podataka iz forme
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Poziv funkcije za proveru korisnika

		var valid bool = data.ValidUser(username, password)
		if valid {
			fmt.Fprintf(w, "Da")
		} else {
			fmt.Fprintf(w, "Ne")
		}

	})

	// Obrada podataka iz forme kada je ruta /register
	mux.HandleFunc("/register", func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodPost {
			http.Error(res, "Only POST method allowed", http.StatusMethodNotAllowed)
			return
		}

		// Čitanje podataka iz forme
		username := req.FormValue("username")
		password := req.FormValue("password")

		// Poziv funkcije za čuvanje korisnika
		data.SaveUserApplication(username, password)
	})

	// Listen and serve on port 3000
	fmt.Println("[ UP ON PORT 3000 ]")

	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}

func loggedinHandler(w http.ResponseWriter, r *http.Request, githubData string) {
	if githubData == "" {
		// Unauthorized users get an unauthorized message
		fmt.Fprintf(w, "UNAUTHORIZED!")
		return
	}

	w.Header().Set("Content-type", "application/json")

	var prettyJSON bytes.Buffer
	parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	if parserr != nil {
		log.Panic("JSON parse error")
	}

	fmt.Fprintf(w, string(prettyJSON.Bytes()))
}

func githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	githubClientID := getGithubClientID()

	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", githubClientID, "http://localhost:3000/login/github/callback")

	http.Redirect(w, r, redirectURL, 301)
}

func githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	githubAccessToken := getGithubAccessToken(code)

	githubData := getGithubData(githubAccessToken)

	loggedinHandler(w, r, githubData)
}

// Function that retrieves users data from github profile
func getGithubData(accessToken string) string {
	req, reqerr := http.NewRequest("GET", "https://api.github.com/user", nil)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	respbody, _ := ioutil.ReadAll(resp.Body)

	return string(respbody)
}

// Function that calls Github api for access token generation
func getGithubAccessToken(code string) string {

	clientID := getGithubClientID()
	clientSecret := getGithubClientSecret()

	requestBodyMap := map[string]string{"client_id": clientID, "client_secret": clientSecret, "code": code}
	requestJSON, _ := json.Marshal(requestBodyMap)

	req, reqerr := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(requestJSON))
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	respbody, _ := ioutil.ReadAll(resp.Body)

	// Represents the response received from Github
	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	var ghresp githubAccessTokenResponse
	json.Unmarshal(respbody, &ghresp)

	return ghresp.AccessToken
}

// Function that gets client id from env file
func getGithubClientID() string {

	githubClientID, exists := os.LookupEnv("CLIENT_ID")
	if !exists {
		log.Fatal("Github Client ID not defined in .env file")
	}

	return githubClientID
}

// Function that gets client secret from env file
func getGithubClientSecret() string {

	githubClientSecret, exists := os.LookupEnv("CLIENT_SECRET")
	if !exists {
		log.Fatal("Github Client ID not defined in .env file")
	}

	return githubClientSecret
}
