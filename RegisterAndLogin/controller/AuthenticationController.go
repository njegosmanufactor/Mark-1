package controller

import (
	user "RegisterAndLogin/model"
	"encoding/json"
	"net/http"
)

type LoginDto struct {
	Username string
	Password string
}

var users = []user.User{
	{Id: 1, Username: "Ranita", Password: "sifra123"},
}

// Implements basic log in functionality with simple authorisation
func LogIn(w http.ResponseWriter, r *http.Request) {
	var logDto LoginDto
	err := json.NewDecoder(r.Body).Decode(&logDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if logDto.Username == users[0].Username {
		if logDto.Password == users[0].Password {
			json.NewEncoder(w).Encode(logDto.Username)
		} else {
			json.NewEncoder(w).Encode("Wrong password")
		}
	} else {
		json.NewEncoder(w).Encode("Wrong username")
	}
}

// Implements basic registration. Saves registered user localy in array users
func Register(w http.ResponseWriter, r *http.Request) {
	var user user.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	users = append(users, user)
	json.NewEncoder(w).Encode(users)
}
