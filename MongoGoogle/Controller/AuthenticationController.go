package Controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	dataBase "MongoGoogle/Repository"
	applicationService "MongoGoogle/Service/ApplicationService"
	companyService "MongoGoogle/Service/CompanyService"
	gitService "MongoGoogle/Service/GitService"
	googleService "MongoGoogle/Service/GoogleService"
	ownerService "MongoGoogle/Service/OwnerService"
	tokenService "MongoGoogle/Service/TokenService"
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

	//Admin or owner sends invitation mail. Body requiers company id and user email.
	r.HandleFunc("/sendInvitation", func(res http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("Authorization")
		token = tokenService.SplitTokenHeder(token)
		_, tokenUser, _ := tokenService.ExtractUserFromToken(token)
		if tokenUser != nil && tokenUser.Valid {
			ownerService.SendInvitation(res, req)
		} else {
			json.NewEncoder(res).Encode("Session timed out or terminated")
		}
	})
	//Link that users clicks in his mail message. Writes company id to his comany field.
	r.HandleFunc("/inviteConfirmation/{id}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		transactionId := vars["id"]
		applicationService.IncludeUserInCompany(transactionId, res)
	})
	//Users request for changing forgotten password
	r.HandleFunc("/forgotPassword", func(res http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("Authorization")
		token = tokenService.SplitTokenHeder(token)
		_, tokenUser, _ := tokenService.ExtractUserFromToken(token)
		if tokenUser != nil && tokenUser.Valid {
			applicationService.PasswordChange(res, req)
		} else {
			json.NewEncoder(res).Encode("Session timed out or terminated")
		}
	})
	// Handles forgotten password update after callback.
	r.HandleFunc("/forgotPassword/callback/{transferId}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		transferId := vars["transferId"]
		applicationService.FinaliseForgottenPasswordUpdate(transferId, res, req)
	})
	//Owner send mail to user which he intends to transfer ownership to. Body has owners id,company id and users email
	r.HandleFunc("/trasferOwnership", func(res http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("Authorization")
		token = tokenService.SplitTokenHeder(token)
		_, tokenUser, err := tokenService.ExtractUserFromToken(token)
		if err != nil {
			http.Error(res, "Error extracting user from token", http.StatusInternalServerError)
			return
		}
		if tokenUser != nil && tokenUser.Valid {
			ownerService.TransferOwnership(res, req)
		} else {
			json.NewEncoder(res).Encode("Session timed out or terminated")
		}
	})
	//Sets users field "Role" to "Owner" DA LI UBACITI DA SE PROSLI OWNER OBRISE?
	r.HandleFunc("/transferOwnership/feedataBaseack/{transferId}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		transferId := vars["transferId"]
		ownerService.FinaliseOwnershipTransfer(transferId, res)
	})
	// Handles Google login logic, completes user authentication, generates a token, and returns it.
	r.HandleFunc("/googleLogin", func(res http.ResponseWriter, req *http.Request) {
		accessToken := req.URL.Query().Get("access_token")
		googleUser := tokenService.TokenGoogleLoginLogic(res, req, accessToken)
		googleService.CompleteGoogleUserAuthentication(res, req, googleUser)
		user, _ := dataBase.GetUserData(googleUser.Email)
		tokenString, _ := tokenService.GenerateToken(user, time.Hour)
		if tokenString != "" {
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
	// Initiates the process of sending a magic link for login without password.
	r.HandleFunc("/magicLink", func(res http.ResponseWriter, req *http.Request) {
		applicationService.MagicLink(res, req)
	})
	// Confirms the magic link for login and generates a token.
	r.HandleFunc("/confirmMagicLink", func(res http.ResponseWriter, req *http.Request) {
		var requestBody struct {
			Email string `json:"email"`
		}
		errReq := json.NewDecoder(req.Body).Decode(&requestBody)
		if errReq != nil {
			http.Error(res, "Error decoding request body", http.StatusBadRequest)
			return
		}
		user, _ := dataBase.GetUserData(requestBody.Email)
		tokenString, _ := tokenService.GenerateToken(user, time.Hour)
		json.NewEncoder(res).Encode(tokenString)
	})
	// Handles the request for a password-less login code.
	r.HandleFunc("/passwordLessCode", func(res http.ResponseWriter, req *http.Request) {
		applicationService.PasswordLessCode(res, req)
	})
	// Confirms the password-less login code and generates a token.
	r.HandleFunc("/passwordLessCodeConfirm", func(res http.ResponseWriter, req *http.Request) {
		var requestBody struct {
			RequestID string `json:"requestID"`
			Code      string `json:"code"`
		}
		errReq := json.NewDecoder(req.Body).Decode(&requestBody)
		if errReq != nil {
			http.Error(res, "Error decoding request body", http.StatusBadRequest)
			return
		}
		result, _ := dataBase.FindCodeRequestByHex(requestBody.RequestID, res)
		if requestBody.Code == result.Code {
			user, _ := dataBase.GetUserData(result.Email)
			token, _ := tokenService.GenerateToken(user, time.Hour)
			dataBase.DeletePandingRequrst(requestBody.RequestID)
			json.NewEncoder(res).Encode(token)
		} else {
			json.NewEncoder(res).Encode("Incorect code")
		}
	})

	//Our service that serves registration functionality
	r.HandleFunc("/register", func(res http.ResponseWriter, req *http.Request) {
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
		json.NewEncoder(res).Encode(tokenExpString)
	})

	// Verifies the user with the specified email address.
	r.HandleFunc("/verify/{email}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		email := vars["email"]
		if dataBase.VerifyUser(email) {
			fmt.Println(res, email)
		}
	})

	// Registers a new company using the provided email address for authentication.
	r.HandleFunc("/registerCompany", func(res http.ResponseWriter, req *http.Request) {
		companyService.CreateComapny(res, req)
	})

	// Deletes the company associated with the provided email address.
	r.HandleFunc("/deleteCompany", func(res http.ResponseWriter, req *http.Request) {
		companyService.DeleteCompany(res, req)
	})

	//Mux router listens for requests on port : 3000
	fmt.Println("[ UP ON PORT 3000 ]")
	err := http.ListenAndServe(":3000", r)
	log.Fatal(err)
}
