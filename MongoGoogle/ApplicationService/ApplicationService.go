package ApplicationService

import (
	conn "MongoGoogle/Repository"
	data "MongoGoogle/Repository"
	"context"

	"fmt"
	"net/http"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// ApplicationRegister validates user input for registration and saves the application if valid, sending a verification email upon success.
func ApplicationRegister(email string, firstName string, lastName string, phone string, date string, username string, password string) {
	if email == "" || username == "" || password == "" || date == "" || phone == "" || firstName == "" || lastName == "" {
		fmt.Println("Some required parameters are missing.")
		return
	}
	dateOfBirth, err := time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Println("Invalid date format.")
		return
	}
	if dateOfBirth.After(time.Now()) {
		fmt.Println("Date of birth cannot be in the future.")
		return
	}
	match, _ := regexp.MatchString("^[0-9]+$", phone)
	if !match {
		fmt.Println("Phone number must contain only digits.")
		return
	}
	//Save user
	fmt.Println(username)
	if data.ValidEmail(email) {
		fmt.Println("Email in use")
		return
	}
	if data.ValidUsername(username) {
		fmt.Println("Username in use")
		return
	} else {
		data.SaveUserApplication(email, firstName, lastName, phone, date, username, password, false)
		SendMail(email)
		fmt.Println("Success")
	}
}

// Authenticates the user by verifying the email and password, and extracts user information from the token in the request header to set the user as authorized.
func ApplicationLogin(email string, password string) string {
	if !data.ValidUser(email, password) {
		return "Incorrect email or password"
	}
	return "Success"
}

// Includes the user in the company by updating the company ID in the user's document.
func IncludeUserInCompany(companyID string, email string, res http.ResponseWriter) error {
	collection := conn.GetClient().Database("UserDatabase").Collection("Users")
	filter := bson.M{"Email": email}
	update := bson.M{"$set": bson.M{"Company": companyID}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
