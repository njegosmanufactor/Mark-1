package ApplicationService

import (
	model "MongoGoogle/Model"
	conn "MongoGoogle/Repository"
	data "MongoGoogle/Repository"
	"context"
	"encoding/json"
	"log"

	"fmt"
	"net/http"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
func IncludeUserInCompany(requestId string, res http.ResponseWriter) {

	//Finding the right pending request
	collection := conn.GetClient().Database("UserDatabase").Collection("PendingRequests")
	requestIdentifier, iderr := primitive.ObjectIDFromHex(requestId)
	if iderr != nil {
		log.Fatal(iderr)
	}
	filter := bson.M{"_id": requestIdentifier}
	var result model.PendingRequest
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			json.NewEncoder(res).Encode("Didnt find request!")
		}
		log.Fatal(err)
	}
	//Inserts user to company employees field
	conn.AddUserToCompany(result.CompanyID, result.Email, res)
	//Updating pending request to completed
	update := bson.M{"$set": bson.M{"Completed": true}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {

		json.NewEncoder(res).Encode("Table not updated!")
		log.Fatal(err)
	}
}
