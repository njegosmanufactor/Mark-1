package ApplicationService

import (
	data "MongoGoogle/MongoDB"

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
	phone := req.FormValue("phone")
	date := req.FormValue("date")
	username := req.FormValue("username")
	password := req.FormValue("password")
	company := req.FormValue("company")
	country := req.FormValue("country")
	city := req.FormValue("city")
	address := req.FormValue("address")

	//Save user
	data.SaveUserApplication(email, firstName, lastName, phone, date, username, password, company, country, city, address)

	//morace kojic da doda verified polje u model korisnika
	SendMail(email)
}

func ApplicationLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	// Reading from html
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Validation User for Application

	var valid bool = data.ValidUser(username, password)
	if valid {
		fmt.Fprintf(w, "Successful")
	} else {
		fmt.Fprintf(w, "Incorrect username or password")
	}
}
