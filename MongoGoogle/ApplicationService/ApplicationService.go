package ApplicationService

import (
	model "MongoGoogle/Model"
	conn "MongoGoogle/Repository"
	data "MongoGoogle/Repository"
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"

	"fmt"
	"net/http"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		data.SaveUserApplication(email, firstName, lastName, phone, date, username, password, false, "Application")
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
func IncludeUserInCompany(requestId string, res http.ResponseWriter) {

	collection := conn.GetClient().Database("UserDatabase").Collection("PendingRequests")
	requestIdentifier, iderr := primitive.ObjectIDFromHex(requestId)
	if iderr != nil {
		log.Fatal(iderr)
	}
	filter := bson.M{"_id": requestIdentifier}
	var result model.PendingRequest
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(res).Encode("Didnt find user!")
		}
	}
	//ubaciti usera u listu zaposlenih u kompaniji
	conn.AddUserToCompany(result.CompanyID, result.Email, res)
	//Azurirati completed na true
	update := bson.M{"$set": bson.M{"Completed": true}}
	// Perform the update operation
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
		json.NewEncoder(res).Encode("Table not updated!")
	}
}
