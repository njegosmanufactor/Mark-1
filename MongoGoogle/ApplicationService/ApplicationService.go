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

func IncludeUserInCompany(companyID string, email string, res http.ResponseWriter) error {
	collection := conn.GetClient().Database("UserDatabase").Collection("Users")
	filter := bson.M{"Email": email}
	//Ovde mozda bude moralo da se parsira id firme na ObjectID("blablabla")
	update := bson.M{"$set": bson.M{"Company": companyID}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
