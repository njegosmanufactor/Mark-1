package Controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	service "MongoGoogle/Service"
)

type AuthenticationController struct {
	Router *mux.Router
}

func NewAuthenticationController() *AuthenticationController {
	return &AuthenticationController{
		Router: mux.NewRouter(),
	}
}

func (ac *AuthenticationController) RegisterRoutes() {
	ac.Router.HandleFunc("/auth/login/github", ac.GithubLogin)
	ac.Router.HandleFunc("/auth/login/github/callback", ac.GithubLoginCallback)
	ac.Router.HandleFunc("/auth/googleLogin", ac.GoogleLogIn)
	ac.Router.HandleFunc("/auth/login", ac.Login)
	ac.Router.HandleFunc("/auth/magicLing", ac.MagicLink)
	ac.Router.HandleFunc("/auth/confirmMagicLink/{email}", ac.MagicLinkCallback)
	ac.Router.HandleFunc("/auth/passwordLessCode", ac.PasswordlessLogin)
	ac.Router.HandleFunc("/auth/passwordLessCodeConfirm", ac.PasswordlessLoginCallback)
	ac.Router.HandleFunc("/auth/register", ac.Register)
	ac.Router.HandleFunc("/auth/logout", ac.Logout)
	ac.Router.HandleFunc("/auth/verify/{email}", ac.VerifyEmail)
}

func (ac *AuthenticationController) GithubLogin(res http.ResponseWriter, req *http.Request) {
	service.GithubLoginHandler(res, req)
}
func (ac *AuthenticationController) GithubLoginCallback(res http.ResponseWriter, req *http.Request) {
	service.GithubCallbackHandler(res, req)
}
func (ac *AuthenticationController) GoogleLogIn(res http.ResponseWriter, req *http.Request) {
	accessToken := req.URL.Query().Get("access_token")
	googleUser := service.TokenGoogleLoginLogic(res, req, accessToken)
	service.CompleteGoogleUserAuthentication(res, req, googleUser)
	user := service.GetUserData(googleUser.Email)
	tokenString, _ := service.GenerateToken(user, time.Hour)
	if tokenString != "" {
		json.NewEncoder(res).Encode(tokenString)
	}
}
func (ac *AuthenticationController) Login(res http.ResponseWriter, req *http.Request) {
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
	service.TokenAppLoginLogic(res, req, authHeader, requestBody.Email, requestBody.Password)
}
func (ac *AuthenticationController) MagicLink(res http.ResponseWriter, req *http.Request) {
	service.MagicLink(res, req)
}
func (ac *AuthenticationController) MagicLinkCallback(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	email := vars["email"]
	user := service.GetUserData(email)
	tokenString, _ := service.GenerateToken(user, time.Hour)
	json.NewEncoder(res).Encode(tokenString)
}
func (ac *AuthenticationController) PasswordlessLogin(res http.ResponseWriter, req *http.Request) {
	service.PasswordLessCode(res, req)
}

func (ac *AuthenticationController) PasswordlessLoginCallback(res http.ResponseWriter, req *http.Request) {
	var requestBody struct {
		RequestID string `json:"requestID"`
		Code      string `json:"code"`
	}
	errReq := json.NewDecoder(req.Body).Decode(&requestBody)
	if errReq != nil {
		http.Error(res, "Error decoding request body", http.StatusBadRequest)
		return
	}
	result, _ := service.FindCodeRequestByHex(requestBody.RequestID, res)
	if requestBody.Code == result.Code {
		user := service.GetUserData(result.Email)
		token, _ := service.GenerateToken(user, time.Hour)
		service.DeletePendingRequest(requestBody.RequestID)
		json.NewEncoder(res).Encode(token)
	} else {
		json.NewEncoder(res).Encode("Incorect code")
	}
}

func (ac *AuthenticationController) Register(res http.ResponseWriter, req *http.Request) {
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
	service.ApplicationRegister(requestBody.Email, requestBody.FirstName, requestBody.LastName, requestBody.PhoneNumber, requestBody.Date, requestBody.Username, requestBody.Password, res)

}
func (ac *AuthenticationController) Logout(res http.ResponseWriter, req *http.Request) {
	user, _ := service.GetUserAndPointerFromToken(res, req)
	tokenExpString, _ := service.GenerateToken(user, time.Second)
	json.NewEncoder(res).Encode(tokenExpString)
}

func (ac *AuthenticationController) VerifyEmail(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	email := vars["email"]
	if service.VerifyUser(email) {
		fmt.Println(res, email)
	}
}
