package ApplicationService

import (
	model "MongoGoogle/Model"
	conn "MongoGoogle/Repository"
	"context"
	"encoding/json"
	"log"
	"unicode"

	"fmt"
	"net/http"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PasswordChangeDTO struct {
	Email string `bson:"email"`
}

type NewPassword struct {
	Password        string `bson:"Password"`
	ConfirmPassword string `bson:"ConfirmedPassword"`
}

// Function that checks if the string contains given symbols.
func containsSpecialCharacters(input string) bool {
	pattern := regexp.MustCompile(`[!@#$%^&*()?><,./|\}{=-_+]`)
	if match := pattern.FindStringIndex(input); match != nil {
		return true
	}

	return false
}

// Treba da se kreira pending zahtev u kom ce da stoji completed i email korisnika koji zeli da promeni sifru.
// Nakon kreiranja te tabele korisniku se u mejl salje link sa _id tog zahteva, da ne bi mogao da se unsese email preko putanje
// Na taj nacin se onemogucuje da se promeni sifra bilo kog korisnika navodjenjem mejla
func PasswordChange(res http.ResponseWriter, req *http.Request) {
	var passChangeDTO PasswordChangeDTO
	decErr := json.NewDecoder(req.Body).Decode(&passChangeDTO)
	if decErr != nil {
		http.Error(res, decErr.Error(), http.StatusBadRequest)
	}
	user, found := conn.FindUserByMail(passChangeDTO.Email, res)
	if found {
		if user.Verified {
			_, id := conn.CreatePasswordChangeRequest(passChangeDTO.Email)
			SendPasswordChangeLink(id.Hex(), passChangeDTO.Email)
		} else {
			json.NewEncoder(res).Encode("This user hasn't verified his account.")
		}
	} else {
		json.NewEncoder(res).Encode("Didnt find the user!")
	}
}
func FinaliseForgottenPasswordUpdate(transferId string, res http.ResponseWriter, req *http.Request) {
	var password NewPassword
	decErr := json.NewDecoder(req.Body).Decode(&password)
	if decErr != nil {
		http.Error(res, decErr.Error(), http.StatusBadRequest)
	}
	if password.Password != password.ConfirmPassword {
		json.NewEncoder(res).Encode("Passwords dont match!")
		return
	}
	collection := conn.GetClient().Database("UserDatabase").Collection("PendingRequests")
	userCollection := conn.GetClient().Database("UserDatabase").Collection("Users")
	requestIdentifier, iderr := primitive.ObjectIDFromHex(transferId)
	if iderr != nil {
		log.Fatal(iderr)
	}
	filter := bson.M{"_id": requestIdentifier}
	var result model.PasswordChangeRequest
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			json.NewEncoder(res).Encode("Didnt find request!")
		}
		log.Fatal(err)
	}
	//Updating the password field in user
	NewPassword, _ := conn.HashPassword(password.Password)
	userUpdate := bson.M{"$set": bson.M{"Password": NewPassword}}
	userFilter := bson.M{"Email": result.Email}
	_, userErr := userCollection.UpdateOne(context.Background(), userFilter, userUpdate)
	if userErr != nil {
		json.NewEncoder(res).Encode("Password not updated!")
	}

	//Updating the request table
	update := bson.M{"$set": bson.M{"Completed": true}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		json.NewEncoder(res).Encode("Table not updated!")
	}
}

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
	if len(password) < 6 {
		fmt.Println("Password must contain at least 6 charracters.")
		return
	}
	if !containsSpecialCharacters(password) {
		fmt.Println("Password must contain a special charracter!")
		return
	}
	HasUpper := false
	HasLower := true
	for _, r := range password {
		if unicode.IsUpper(r) {
			HasUpper = true
		}
		if unicode.IsLower(r) {
			HasLower = true
		}
	}
	if !(HasUpper && HasLower) {
		fmt.Println("Password must contain uppercase and lowercase letters!")
		return
	}
	//Save user
	if conn.FindUserEmail(email) {
		fmt.Println("Email in use")
		return
	}
	if conn.FindUserUsername(username) {
		fmt.Println("Username in use")
		return
	} else {
		hashedPass, hashError := conn.HashPassword(password)
		if hashError != nil {
			log.Panic(hashError)
		}
		conn.SaveUserApplication(email, firstName, lastName, phone, date, username, hashedPass, false, "Application")
		SendMail(email)
		fmt.Println("Success")
	}
}

// Authenticates the user by verifying the email and password, and extracts user information from the token in the request header to set the user as authorized.
func ApplicationLogin(email string, password string) string {
	fmt.Println(password)
	if !conn.ValidUser(email, password) {
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
