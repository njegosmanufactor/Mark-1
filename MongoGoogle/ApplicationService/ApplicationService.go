package ApplicationService

import (
	data "MongoGoogle/Repository"

	"fmt"
	"net/http"
)

func ApplicationRegister(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		http.Error(res, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	email := req.FormValue("email")
	firstName := req.FormValue("firstName")
	lastName := req.FormValue("lastName")
	phoneNumber := req.FormValue("countryCode") + req.FormValue("phone")
	date := req.FormValue("date")
	username := req.FormValue("username")
	password := req.FormValue("password")

	//Save user
	if data.ValidEmail(email) || data.ValidUsername(username) {
		fmt.Fprintf(res, "Username or Email in use")
	} else {
		data.SaveUserApplication(email, firstName, lastName, phoneNumber, date, username, password)
	}
	SendMail(email)
	http.Redirect(res, req, "success.html", http.StatusSeeOther)
}

func ApplicationLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	// Čitanje podataka iz forme
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Provera korisničkih podataka
	if !data.ValidUser(email, password) {
		http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
		return
	}

	// Preusmeravanje na success.html ako je prijava uspešna
	http.Redirect(w, r, "/success.html", http.StatusSeeOther)
}
