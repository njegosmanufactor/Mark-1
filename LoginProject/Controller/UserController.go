package Controller

import (
	model "LoginProject/Model"
	"encoding/json"
	"fmt"
	"net/http"
)

// Globalna lista korisnika
var users []model.User

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users = append(users, user)

	fmt.Fprintf(w, "Korisnik uspješno registrovan: %+v\n", user)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Vrati sve korisnike kao JSON
	json.NewEncoder(w).Encode(users)
}

// Prijavljivanje korisnika
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
			fmt.Fprintf(w, "Korisnik %s je uspješno prijavljen.\n", user.Username)
			return
		}
	}

	http.Error(w, "Pogrešno korisničko ime ili šifra.", http.StatusUnauthorized)
}

// Promjena šifre korisnika
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
			// Pronađen korisnik, mijenjamo šifru
			users[i].Password = changePasswordInfo.NewPassword
			fmt.Fprintf(w, "Šifra za korisnika %s uspješno promenjena.\n", user.Username)
			return
		}
	}

	http.Error(w, "Pogrešno korisničko ime ili šifra.", http.StatusUnauthorized)
}
