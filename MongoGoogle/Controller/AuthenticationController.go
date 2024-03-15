package Controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"

	applicationService "MongoGoogle/ApplicationService"
	gitService "MongoGoogle/GitService"
	googleService "MongoGoogle/GoogleService"
	userType "MongoGoogle/Model"
	ownerService "MongoGoogle/OwnerService"
	db "MongoGoogle/Repository"
	tokenService "MongoGoogle/TokenService"
)

// Handles user authentication using OAuth2 providers such as Google and Github.
func Mark1() {
	//Client secret created on google cloud platform/ Apis & Services / Credentials
	var key, env_key_error = os.LookupEnv("GOOGLE_KEY")
	if !env_key_error {
		log.Fatal("Google key not defined in .env file")
	}
	var client_id, env_clientID_error = os.LookupEnv("GOOGLE_CLIENT_ID")
	if !env_clientID_error {
		log.Fatal("Google client id not defined in .env file")
	}
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
		google.New(client_id, key, "http://localhost:3000/auth/google/callback", "email", "profile"),
	)

	r := mux.NewRouter()

	//Using google OAuth2 to authenticate user
	r.HandleFunc("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	r.HandleFunc("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
		googleService.CompleteGoogleUserAuthentication(res, req)
	})

	//Github authentication paths
	r.HandleFunc("/login/github", func(w http.ResponseWriter, r *http.Request) {
		gitService.GithubLoginHandler(w, r)
	})

	// Github callback
	r.HandleFunc("/login/github/callback", func(w http.ResponseWriter, r *http.Request) {
		gitService.GithubCallbackHandler(w, r)
	})

	// Route where the authenticated user is redirected to
	r.HandleFunc("/loggedin", func(w http.ResponseWriter, r *http.Request) {
		var nill userType.GitHubData
		gitService.LoggedinHandler(w, r, nill)
	})
	//Funkcija koju admin klikce, znaci treba da se u njenom body nalaze mail korisnika, i id kompanije.
	r.HandleFunc("/sendInvitation", func(res http.ResponseWriter, req *http.Request) { //postman
		ownerService.SendInvitation(res, req)
	})
	//funkcija koja ce da upisuje id kompanije u korisnikov profil u bazi
	r.HandleFunc("/inviteConfirmation/{companyID}/{userID}", func(res http.ResponseWriter, req *http.Request) { //postman
		vars := mux.Vars(req)
		companyID := vars["companyID"]
		userID := vars["userID"]
		applicationService.IncludeUserInCompany(companyID, userID, res)
	})
	r.HandleFunc("/trasferOwnership", func(res http.ResponseWriter, req *http.Request) { //postman
		ownerService.TransferOwnership(res, req)
	})
	r.HandleFunc("/transferOwnership/feedback/{email}", func(res http.ResponseWriter, req *http.Request) { //postman
		vars := mux.Vars(req)
		email := vars["email"]
		ownerService.FinaliseOwnershipTransfer(email)
	})

	//Our service that serves login functionality
	r.HandleFunc("/login", func(res http.ResponseWriter, req *http.Request) {
		var requestBody struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(req.Body).Decode(&requestBody)
		if err != nil {
			http.Error(res, "Error decoding request body", http.StatusBadRequest)
			return
		}
		authHeader := req.Header.Get("Authorization")

		//ako nema token
		if authHeader == "" {
			user, _ := db.GetUserData(requestBody.Email)
			token, _ := tokenService.GenerateToken(user)
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode(struct {
				Token string `json:"token"`
			}{
				Token: token,
			})
		} else { //ako ima token
			authHeader = tokenService.SplitTokenHeder(authHeader)
			user, token, _ := tokenService.ExtractUserFromToken(authHeader)
			if err != nil {
				http.Error(res, "Error extracting user from token", http.StatusInternalServerError)
				return
			}

			if token != nil && token.Valid {
				// Token nije nil i važeći je
				message := applicationService.ApplicationLogin(requestBody.Email, requestBody.Password)
				res.Header().Set("Content-Type", "application/json")
				json.NewEncoder(res).Encode(struct {
					Message     string                   `json:"message"`
					Token       *jwt.Token               `json:"token"`
					CurrentUser userType.ApplicationUser `json:"user"`
				}{
					Message:     message,
					Token:       token,
					CurrentUser: user,
				})
			} else {
				// Token je ili nil ili nevažeći
				user, _ := db.GetUserData(requestBody.Email)
				newToken, _ := tokenService.GenerateToken(user)
				res.Header().Set("Content-Type", "application/json")
				json.NewEncoder(res).Encode(struct {
					Token string `json:"token"`
				}{
					Token: newToken,
				})
			}

		}

	})

	//Our service that serves registration functionality
	r.HandleFunc("/register", func(res http.ResponseWriter, req *http.Request) { //postman
		var requestBody struct {
			Email       string `json:"email"`
			FirstName   string `json:"firstName"`
			LastName    string `json:"lastName"`
			PhoneNumber string `json:"phoneNumber"`
			Date        string `json:"date"`
			Username    string `json:"username"`
			Password    string `json:"password"`
		}
		err := json.NewDecoder(req.Body).Decode(&requestBody)
		if err != nil {
			http.Error(res, "Error decoding request body", http.StatusBadRequest)
			return
		}
		applicationService.ApplicationRegister(requestBody.Email, requestBody.FirstName, requestBody.LastName, requestBody.PhoneNumber, requestBody.Date, requestBody.Username, requestBody.Password)
	})

	//NE RADI
	// Logs out the user with the specified email address.
	r.HandleFunc("/logout", func(res http.ResponseWriter, req *http.Request) {
		tokenString := req.Header.Get("Authorization")
		tokenString = tokenService.SplitTokenHeder(tokenString)
		user, token, err := tokenService.ExtractUserFromToken(tokenString)
		if err != nil {
			http.Error(res, "Error extracting user from token", http.StatusInternalServerError)
			return
		}
		tokenService.SetTokenExpired(token)
		fmt.Println("Logout" + " " + user.Email)
		userLogout, _ := db.GetUserData(user.Email)
		newToken, _ := tokenService.GenerateToken(userLogout)
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(struct {
			Token string `json:"token"`
		}{
			Token: newToken,
		})
	})

	// Verifies the user with the specified email address.
	r.HandleFunc("/verify/{email}", func(res http.ResponseWriter, req *http.Request) { //postman

		vars := mux.Vars(req)
		email := vars["email"]
		if db.VerifyUser(email) {
			fmt.Fprintf(res, email)
		}

	})

	/////////////////////////////////  COMPANY    ///////////////////////////////////

	// Registers a new company using the provided email address for authentication.
	r.HandleFunc("/registerCompany", func(res http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("Authorization")
		token = tokenService.SplitTokenHeder(token)
		user, tokenUser, err := tokenService.ExtractUserFromToken(token)
		if err != nil {
			http.Error(res, "Error extracting user from token", http.StatusInternalServerError)
			return
		}
		if tokenUser != nil && tokenUser.Valid {
			var companyData struct {
				Name                  string            `json:"name"`
				Address               userType.Location `json:"location"`
				Website               string            `json:"website"`
				ListOfApprovedDomains []string          `json:"listOfApprovedDomains"`
			}

			err := json.NewDecoder(req.Body).Decode(&companyData)
			if err != nil {
				http.Error(res, "Error decoding request body", http.StatusBadRequest)
				return
			}

			if db.ValidComapnyName(companyData.Name) {
				fmt.Printf("Company exist\n")
			} else {
				db.SetUserRole(user.ID, "Owner")
				db.SetOwnerCompany(companyData.Name, user.ID.String())
				db.SaveCompany(companyData.Name, companyData.Address, companyData.Website, companyData.ListOfApprovedDomains, user.ID)
				user, _ := db.GetUserData(user.Email)
				token, _ := tokenService.GenerateToken(user)
				res.Header().Set("Content-Type", "application/json")
				json.NewEncoder(res).Encode(struct {
					Token string `json:"token"`
				}{
					Token: token,
				})
			}
		} else {
			fmt.Println("User not found")
		}

	})

	// Deletes the company associated with the provided email address.
	r.HandleFunc("/deleteCompany", func(res http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("Authorization")
		token = tokenService.SplitTokenHeder(token)
		user, tokenUser, err := tokenService.ExtractUserFromToken(token)
		if err != nil {
			http.Error(res, "Error extracting user from token", http.StatusInternalServerError)
			return
		}

		var requestBody struct {
			CompanyName string `json:"companyName"`
		}
		errReq := json.NewDecoder(req.Body).Decode(&requestBody)
		if errReq != nil {
			http.Error(res, "Error decoding request body", http.StatusBadRequest)
			return
		}
		if tokenUser != nil && tokenUser.Valid {
			if requestBody.CompanyName == "" {
				http.Error(res, "Company name is required", http.StatusBadRequest)
				return
			}
			db.SetUserRole(user.ID, "User")
			db.DeleteCompany(requestBody.CompanyName)
			user, _ := db.GetUserData(user.Email)
			token, _ := tokenService.GenerateToken(user)
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode(struct {
				Token string `json:"token"`
			}{
				Token: token,
			})
		} else {
			fmt.Println("You are not owner of" + requestBody.CompanyName)
		}
	})

	//Mux router listens for requests on port : 3000
	fmt.Println("[ UP ON PORT 3000 ]")
	err := http.ListenAndServe(":3000", r)
	log.Fatal(err)
}
