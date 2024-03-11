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
	phoneNumber := req.FormValue("countryCode") + req.FormValue("phone")
	date := req.FormValue("date")
	username := req.FormValue("username")
	password := req.FormValue("password")
	company := req.FormValue("company")
	country := req.FormValue("country")
	city := req.FormValue("city")
	address := req.FormValue("address")

	//Save user
	if data.ValidEmail(email) || data.ValidUsername(username) {
		fmt.Fprintf(res, "Username or Email in use")
	} else {
		data.SaveUserApplication(email, firstName, lastName, phoneNumber, date, username, password, company, country, city, address)
	}
	SendMail(email)
}

func ApplicationLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	// Reading from html
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Validation User for Application

	var valid bool = data.ValidUser(email, password)
	if valid {
		fmt.Fprintf(w, "Successful")
		user, err := data.GetUserData(email)
		if err != nil {
			fmt.Printf("User not found: %v\n", err)
			return
		}
		if user.Company != "" {
			data.SetUserRoleOwner(email)
		}
	} else {
		fmt.Fprintf(w, "Incorrect email or password")
	}
}
