package ApplicationService

import (
	conn "MongoGoogle/Repository"
	data "MongoGoogle/Repository"
	"context"
	"encoding/base64"

	"strings"

	"fmt"
	"net/http"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

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
func ApplicationLogin(email string, password string, req *http.Request) {
	if !data.ValidUser(email, password) {
		fmt.Println("Incorrect email or password")
		return
	}
	fmt.Println("Success")
	ExtractUserInfoFromToken(req)
}

// Extracts user information from the token in the request header and sets the user as authorized in the database.
func ExtractUserInfoFromToken(req *http.Request) bool {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		fmt.Println("Unauthorised")
		return false
	}
	parts := strings.Split(authHeader, " ")
	token := parts[1]

	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		fmt.Println("Failed to decode token")
		fmt.Println("Unauthorised")
		return false
	}
	tokenData := strings.Split(string(decodedToken), ":")
	email := tokenData[0]

	user, err := data.GetUserData(email)
	data.SetAuthorise(user.ID, true)
	fmt.Println(email + " " + "Authorized")
	return true
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
