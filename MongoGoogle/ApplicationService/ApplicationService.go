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

	username := req.FormValue("username")
	password := req.FormValue("password")

	//Save user
	data.SaveUserApplication(username, password)
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
