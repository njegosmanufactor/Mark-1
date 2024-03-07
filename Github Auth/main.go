package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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

// Function called before main initialisation that loads env variables for github oauth
func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {
	// Root route
	// Simply returns a link to the login route
	http.HandleFunc("/", rootHandler)

	// Login route
	http.HandleFunc("/login/github/", githubLoginHandler)

	// Github callback
	http.HandleFunc("/login/github/callback", githubCallbackHandler)

	// Route where the authenticated user is redirected to
	http.HandleFunc("/loggedin", func(w http.ResponseWriter, r *http.Request) {
		var nil GitHubData
		loggedinHandler(w, r, nil)
	})

	// Listen and serve on port 3000
	fmt.Println("[ UP ON PORT 3000 ]")
	log.Panic(
		http.ListenAndServe(":3000", nil),
	)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<a href="/login/github/">LOGIN</a>`)
}

func loggedinHandler(w http.ResponseWriter, r *http.Request, githubData GitHubData) {
	if githubData.Id == 0 {
		fmt.Fprintf(w, "UNAUTHORISED!")
		return
	}

	t, _ := template.ParseFiles("pages/success.html")
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

	//After authorisation callback is called
	//loggedinHandler is used to display user data. Redirecting starts here.
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
