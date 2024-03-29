package Service

import (
	userType "MongoGoogle/Model"
	data "MongoGoogle/Repository"
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
	user, found := data.FindUserByMail(githubData.Username, w)
	if found {
		token, _ := GenerateToken(user, time.Hour)
		json.NewEncoder(w).Encode(token)
	} else {
		fmt.Println("Account created git")
		data.SaveUserApplication(githubData.Username, githubData.Name, "", "", "", githubData.Username, "", true, "GitHub")
	}
}

// GithubLoginHandler redirects the user to GitHub's OAuth authorization page.
func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	githubClientID := GetGithubClientID()
	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", githubClientID, "http://localhost:3000/login/github/callback")
	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}

// Handles the callback from GitHub authentication.
func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	//Ovaj "access_token", mora da se preimenuje u "code" da bi se moglo gadjati preko fronta.
	//Jedino ako moze u postmanu nekako da se querry parametar promeni sa access_token na code, ne bi ovde moralo da se menja
	code := r.URL.Query().Get("access_token")
	//Ovo je za generisanje koda kada se ne gadja postman. Iz postmana kod je direktno githubAccesstoken
	//Kada se gadja preko postmana iz postmana se dobija taj kod.
	//githubAccessToken := GetGithubAccessToken(code)
	//githubData := GetGithubData(githubAccessToken)
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
	var res http.ResponseWriter
	githubClientID, exists := os.LookupEnv("CLIENT_ID")
	if !exists {
		json.NewEncoder(res).Encode("Error looking up client_id from .env file")
		return ""
	}
	return githubClientID
}

// Function that gets client secret from env file
func GetGithubClientSecret() string {
	githubClientSecret, exists := os.LookupEnv("CLIENT_SECRET")
	var res http.ResponseWriter
	if !exists {
		json.NewEncoder(res).Encode("Error looking up client_secret from .env file")
		return ""
	}
	return githubClientSecret
}
