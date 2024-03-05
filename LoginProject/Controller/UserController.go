package Controller

import (
	model "LoginProject/Model"
	"encoding/json"
	"fmt"
	"net/http"
)

var users []model.User

// RegisterHandler handles user registration.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users = append(users, user)

	fmt.Fprintf(w, "User successfully registered: %+v\n", user)
}

// GetUsersHandler returns the list of users.
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

// LoginHandler handles user login.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginInfo struct {
		Username string
		Password string
	}
	err := json.NewDecoder(r.Body).Decode(&loginInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, user := range users {
		if user.Username == loginInfo.Username && user.Password == loginInfo.Password {
			fmt.Fprintf(w, "User %s successfully logged in.\n", user.Username)
			return
		}
	}

	http.Error(w, "Incorrect username or password.", http.StatusUnauthorized)
}

// ChangePasswordHandler handles password change for a user.
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var changePasswordInfo struct {
		Username    string
		OldPassword string
		NewPassword string
	}
	err := json.NewDecoder(r.Body).Decode(&changePasswordInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, user := range users {
		if user.Username == changePasswordInfo.Username && user.Password == changePasswordInfo.OldPassword {
			users[i].Password = changePasswordInfo.NewPassword
			fmt.Fprintf(w, "Password for user %s successfully changed.\n", user.Username)
			return
		}
	}

	http.Error(w, "Incorrect username or password.", http.StatusUnauthorized)
}
