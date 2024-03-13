package Repository

import (
	model "MongoGoogle/Model"
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Setting up the URL to connect to the MongoDB server
var Uri, _ = os.LookupEnv("MONGO_URI")
var ClientOptions = options.Client().ApplyURI(Uri)
var Client, Err = mongo.Connect(context.Background(), ClientOptions)

// save user into database
func SaveUserOther(email string) {
	UsersCollection := Client.Database("UserDatabase").Collection("Users")
	if Err != nil {
		log.Fatal(Err)
	}
	// Creating user instance
	user := model.OtherUser{
		Email: email,
	}

	// Adding user to the database
	insertResult, err := UsersCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added new user with ID:", insertResult.InsertedID)
}

func SaveUserApplication(email string, firstName string, lastName string, phone string, date string, username string, password string) {

	UsersCollection := Client.Database("UserDatabase").Collection("Users")
	if Err != nil {
		log.Fatal(Err)
	}
	// Creating user instance
	user := model.ApplicationUser{
		Email:       email,
		FirstName:   firstName,
		LastName:    lastName,
		Phone:       phone,
		DateOfBirth: date,
		Username:    username,
		Password:    password,
		Company:     "",
		Role:        "User",
	}

	// Adding user to the database
	insertResult, err := UsersCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added new user with ID:", insertResult.InsertedID)
}

func ValidUser(email string, password string) bool {

	UsersCollection := Client.Database("UserDatabase").Collection("Users")
	if Err != nil {
		log.Fatal(Err)
	}

	filter := bson.M{"Email": email, "Password": password}
	var result model.ApplicationUser
	Err = UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	if Err != nil {
		if Err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatal(Err)
	}
	return true
}

func ValidEmail(email string) bool {
	UsersCollection := Client.Database("UserDatabase").Collection("Users")
	if Err != nil {
		log.Fatal(Err)
	}
	filter := bson.M{"Email": email}
	var result model.ApplicationUser
	Err = UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	if Err != nil {
		if Err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatal(Err)
	}
	return true
}

func ValidUsername(username string) bool {
	UsersCollection := Client.Database("UserDatabase").Collection("Users")
	if Err != nil {
		log.Fatal(Err)
	}
	filter := bson.M{"Username": username}
	var result model.ApplicationUser
	Err = UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	if Err != nil {
		if Err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatal(Err)
	}
	return true
}

func VerifyUser(email string) bool {
	UsersCollection := Client.Database("UserDatabase").Collection("Users")
	if Err != nil {
		log.Fatal(Err)
	}
	filter := bson.M{"Email": email}
	update := bson.M{"$set": bson.M{"Verified": true}}
	_, Err = UsersCollection.UpdateOne(context.Background(), filter, update)
	return Err == nil
}

func SetUserRoleOwner(userID primitive.ObjectID) error {

	UsersCollection := Client.Database("UserDatabase").Collection("Users")
	if Err != nil {
		log.Fatal(Err)
	}

	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"Role": "Owner"}}

	_, Err = UsersCollection.UpdateOne(context.Background(), filter, update)
	if Err != nil {
		return Err
	}

	return nil
}
func GetUserData(email string) (model.ApplicationUser, error) {
	UsersCollection := Client.Database("UserDatabase").Collection("Users")
	if Err != nil {
		log.Fatal(Err)
	}

	filter := bson.M{"Email": email}

	var result model.ApplicationUser

	Err = UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	if Err != nil {
		return model.ApplicationUser{}, Err
	}

	return result, nil
}
