package LoginRegister

import (
	"bytes"
	"encoding/json"
	"fmt"

	//"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"

	data "MongoGoogle/MongoDB"
)

type GitHubData struct {
	Login             string  `json:"login"`
	Id                float64 `json:"id"`
	NodeId            string  `json:"node_id"`
	AvatarUrl         string  `json:"avatar_url"`
	GravatarId        string  `json:"gravatar_id"`
	Url               string  `json:"url"`
	HtmlUrl           string  `json:"html_url"`
	FollowersUrl      string  `json:"followers_url"`
	FollowingUrl      string  `json:"following_url"`
	GistsUrl          string  `json:"gists_url"`
	StarredUrl        string  `json:"starred_url"`
	SubscriptionsUrl  string  `json:"subscriptions_url"`
	OrganizationsUrl  string  `json:"organizations_url"`
	ReposUrl          string  `json:"repos_url"`
	EventsUrl         string  `json:"events_url"`
	RecievedEventsUrl string  `json:"recieved_events_url"`
	Type              string  `json:"type"`
	Name              string  `json:"name"`
	Company           string  `json:"company"`
	Blog              string  `json:"blog"`
	Location          string  `json:"location"`
	Email             string  `json:"email"`
	Hireable          bool    `json:"hireable"`
	Bio               string  `json:"bio"`
	TwitterUsername   string  `json:"twitter_username"`
	PublicRepos       int     `json:"public_repos"`
	PublicGists       int     `json:"public_gists"`
	Followers         int     `json:"followers"`
}

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

		// Provjera da li korisničko ime postoji u bazi
		if data.ValidUsername(user.Email) {
			// Korisničko ime postoji, prikazujemo poruku
			fmt.Fprintf(res, "Google Account Successfully Logged In")
		} else {
			// Korisničko ime ne postoji, preusmjeravamo korisnika na Google prijavu
			data.SaveUserOther(user.Email)
			t, _ := template.ParseFiles("LoginRegister/pages/success.html")
			t.Execute(res, user)
		}
		// In success.html file or any front display we can select what we want to display.
		// user is instance of goth.User class of which fields can be found in documentation.
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
		var nill GitHubData
		loggedinHandler(w, r, nill)
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
			fmt.Fprintf(w, "Successful") // Poruka uspjeha
		} else {
			// Poruka neuspjeha
			fmt.Fprintf(w, "Incorrect username or password")
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

func loggedinHandler(w http.ResponseWriter, r *http.Request, githubData GitHubData) {
	if githubData.Id == 0 {
		// Unauthorized users get an unauthorized message
		fmt.Fprintf(w, "Unauthorised")
		return
	}

	t, _ := template.ParseFiles("LoginRegister/pages/success.html")
	t.Execute(w, githubData)
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
func getGithubData(accessToken string) GitHubData {
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

	var data GitHubData
	err := json.Unmarshal(respbody, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}
	return data
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
