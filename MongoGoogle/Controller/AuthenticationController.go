package Controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	applicationService "MongoGoogle/ApplicationService"
	gitService "MongoGoogle/GitService"
	googleService "MongoGoogle/GoogleService"
	userType "MongoGoogle/Model"
	ownerService "MongoGoogle/OwnerService"
	db "MongoGoogle/Repository"
	tokenService "MongoGoogle/TokenService"
)

func Mark1() {

	r := mux.NewRouter()

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

	//Admin or owner sends invitation mail. Body requiers company id and user email.
	r.HandleFunc("/sendInvitation", func(res http.ResponseWriter, req *http.Request) { //postman
		ownerService.SendInvitation(res, req)
	})
	//Link that users clicks in his mail message. Writes company id to his comany field.
	r.HandleFunc("/inviteConfirmation/{id}", func(res http.ResponseWriter, req *http.Request) { //postman
		vars := mux.Vars(req)
		transactionId := vars["id"]
		applicationService.IncludeUserInCompany(transactionId, res)
	})
	//Users request for changing forgotten password
	r.HandleFunc("/forgotPassword", func(res http.ResponseWriter, req *http.Request) {
		applicationService.PasswordChange(res, req)
	})
	r.HandleFunc("/forgotPassword/callback/{transferId}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		transferId := vars["transferId"]
		applicationService.FinaliseForgottenPasswordUpdate(transferId, res, req)
	})
	//Owner send mail to user which he intends to transfer ownership to. Body has owners id,company id and users email
	r.HandleFunc("/trasferOwnership", func(res http.ResponseWriter, req *http.Request) { //postman
		ownerService.TransferOwnership(res, req)
	})
	//Sets users field "Role" to "Owner" DA LI UBACITI DA SE PROSLI OWNER OBRISE?
	r.HandleFunc("/transferOwnership/feedback/{transferId}", func(res http.ResponseWriter, req *http.Request) { //postman
		vars := mux.Vars(req)
		transferId := vars["transferId"]
		ownerService.FinaliseOwnershipTransfer(transferId, res)
	})

	r.HandleFunc("/googleLogin", func(res http.ResponseWriter, req *http.Request) {
		accessToken := req.URL.Query().Get("access_token")
		googleUser := tokenService.TokenGoogleLoginLogic(res, req, accessToken)
		googleService.CompleteGoogleUserAuthentication(res, req, googleUser)
		user, _ := db.GetUserData(googleUser.Email)
		tokenString, _ := tokenService.GenerateToken(user, time.Hour)
		if tokenString != "" {
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode(tokenString)
		}
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
		tokenService.TokenAppLoginLogic(res, req, authHeader, requestBody.Email, requestBody.Password)
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

	// Logs out the user with the specified email address.
	r.HandleFunc("/logout", func(res http.ResponseWriter, req *http.Request) {
		tokenString := req.Header.Get("Authorization")
		tokenString = tokenService.SplitTokenHeder(tokenString)
		user, _, err := tokenService.ExtractUserFromToken(tokenString)
		if err != nil {
			http.Error(res, "Error extracting user from token", http.StatusInternalServerError)
			return
		}
		tokenExpString, _ := tokenService.GenerateToken(user, time.Second)

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode("Token:" + user.Email + "Successfully logged out, this is your new bearer token: " + tokenExpString)
	})

	// Verifies the user with the specified email address.
	r.HandleFunc("/verify/{email}", func(res http.ResponseWriter, req *http.Request) { //postman
		vars := mux.Vars(req)
		email := vars["email"]
		if db.VerifyUser(email) {
			fmt.Println(res, email)
		}
	})

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

			if db.FindComapnyName(companyData.Name) {
				fmt.Printf("Company exist\n")
			} else {
				db.SetUserRole(user.ID, "Owner")
				db.SetOwnerCompany(companyData.Name, user.ID.String())
				db.SaveCompany(companyData.Name, companyData.Address, companyData.Website, companyData.ListOfApprovedDomains, user.ID)
				user, _ := db.GetUserData(user.Email)
				token, _ := tokenService.GenerateToken(user, time.Hour)
				res.Header().Set("Content-Type", "application/json")
				json.NewEncoder(res).Encode("You are successfully created new company: " + companyData.Name + ", this is your new bearer token: " + token)
			}
		} else {
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode("User not found")
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
		if requestBody.CompanyName == "" {
			http.Error(res, "Company name is required", http.StatusBadRequest)
			return
		} else {
			company, err := db.FindCompanyByName(requestBody.CompanyName, res)
			if err == false {
				res.Header().Set("Content-Type", "application/json")
				json.NewEncoder(res).Encode("You don't have any company")
				return
			}
			if tokenUser != nil && tokenUser.Valid && user.ID == company.Owner {
				db.SetUserRole(user.ID, "User")
				db.DeleteCompany(requestBody.CompanyName)
				user, _ := db.GetUserData(user.Email)
				token, _ := tokenService.GenerateToken(user, time.Hour)
				res.Header().Set("Content-Type", "application/json")
				json.NewEncoder(res).Encode("You are successfully deleted company " + requestBody.CompanyName + ", this is your new bearer token: " + token)
			} else {
				res.Header().Set("Content-Type", "application/json")
				json.NewEncoder(res).Encode("You are not owner of" + requestBody.CompanyName)
			}

		}
	})

	//Mux router listens for requests on port : 3000
	fmt.Println("[ UP ON PORT 3000 ]")
	err := http.ListenAndServe(":3000", r)
	log.Fatal(err)
}
