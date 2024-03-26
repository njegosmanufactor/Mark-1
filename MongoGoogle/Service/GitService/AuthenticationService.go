package GitService

import (
	userType "MongoGoogle/Model"
	data "MongoGoogle/Repository"
	tokenService "MongoGoogle/Service/TokenService"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// LoggedinHandler validates if a GitHub user is registered in the database, and if not, it saves the user's GitHub username.
func LoggedinHandler(w http.ResponseWriter, r *http.Request, githubData userType.GitHubData) {
	// Validate user Username in database
	if data.FindUserEmail(githubData.Username) {
		//If user have account
		//Ovde vracaj bearer token
		user, _ := data.FindUserByMail(githubData.Username, w)
		token, _ := tokenService.GenerateToken(user, time.Hour)
		json.NewEncoder(w).Encode(githubData)
		json.NewEncoder(w).Encode(token)

	} else {
		//ovde se registruje
		fmt.Println("Account created git")
		data.SaveUserApplication(githubData.Username, githubData.Name, "", "", "", githubData.Username, "", true, "GitHub")
	}
}

// GithubLoginHandler redirects the user to GitHub's OAuth authorization page.// GithubLoginHandler redirects the user to GitHub's OAuth authorization page.
func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	githubClientID := GetGithubClientID()

	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", githubClientID, "http://localhost:3000/login/github/callback")
	fmt.Println("redirectURL")
	fmt.Println(redirectURL)
	http.Redirect(w, r, redirectURL, 301)

}

// Handles the callback from GitHub authentication.
func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("access_token")
	githubData := GetGithubData(code)
	LoggedinHandler(w, r, githubData)
}

// Function that retrieves users data from github profile
func GetGithubData(accessToken string) userType.GitHubData {
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

	var data userType.GitHubData
	err := json.Unmarshal(respbody, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}
	if data.Role == "" {
		data.Role = "User"
	}
	return data
}

// Function that calls Github api for access token generation
func GetGithubAccessToken(code string) string {

	clientID := GetGithubClientID()
	clientSecret := GetGithubClientSecret()

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
func GetGithubClientID() string {

	githubClientID, exists := os.LookupEnv("CLIENT_ID")
	if !exists {
		log.Fatal("Github Client ID not defined in .env file")
	}

	return githubClientID
}

// Function that gets client secret from env file
func GetGithubClientSecret() string {

	githubClientSecret, exists := os.LookupEnv("CLIENT_SECRET")
	if !exists {
		log.Fatal("Github Client ID not defined in .env file")
	}

	return githubClientSecret
}