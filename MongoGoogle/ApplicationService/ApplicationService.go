package ApplicationService

import (
	conn "MongoGoogle/Repository"
	data "MongoGoogle/Repository"
	"context"

	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
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
			data.SetUserRoleOwner(email) //move to future function
		}
	} else {
		fmt.Fprintf(w, "Incorrect email or password")
	}
}

func IncludeUserInCompany(companyID string, email string, res http.ResponseWriter) error {
	collection := conn.Client.Database("UserDatabase").Collection("Users")
	filter := bson.M{"Email": email}
	//Ovde mozda bude moralo da se parsira id firme na ObjectID("blablabla")
	update := bson.M{"$set": bson.M{"Company": companyID}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
