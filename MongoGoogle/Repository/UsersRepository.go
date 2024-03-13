package Repository

import (
	model "MongoGoogle/Model"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Setting up the URL to connect to the MongoDB server
var Uri = "mongodb+srv://Nikola045:Bombarder535@userdatabase.qcrmscd.mongodb.net/?retryWrites=true&w=majority&appName=UserDataBase"
var ClientOptions = options.Client().ApplyURI(uri)
var Client, Err = mongo.Connect(context.Background(), ClientOptions)

// SaveUserApplication saves user application data into the database.
func SaveUserApplication(email string, firstName string, lastName string, phone string, date string, username string, password string, verified bool) {
	client, err := MongoConnection()
	UsersCollection := client.Database("UserDatabase").Collection("Users")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

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
		Verified:    verified,
		Authorised:  false,
	}

	// Adding user to the database
	insertResult, err := UsersCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added new user with ID:", insertResult.InsertedID)
}

// ValidUser checks if the user with the given email and password exists in the database.
func ValidUser(email string, password string) bool {
	client, err := MongoConnection()
	UsersCollection := client.Database("UserDatabase").Collection("Users")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	filter := bson.M{"Email": email, "Password": password}
	var result model.ApplicationUser
	err = UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatal(err)
	}
	return true
}

// ValidEmail checks if the given email exists in the database.
func ValidEmail(email string) bool {
	client, err := MongoConnection()
	UsersCollection := client.Database("UserDatabase").Collection("Users")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	filter := bson.M{"Email": email}
	var result model.ApplicationUser
	err = UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatal(err)
	}
	return true
}

// ValidUsername checks if the given username exists in the database.
func ValidUsername(username string) bool {
	client, err := MongoConnection()
	UsersCollection := client.Database("UserDatabase").Collection("Users")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	filter := bson.M{"Username": username}
	var result model.ApplicationUser
	err = UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatal(err)
	}
	return true
}

// VerifyUser updates the verification status of the user with the given email.
func VerifyUser(email string) bool {
	client, err := MongoConnection()
	UsersCollection := client.Database("UserDatabase").Collection("Users")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	filter := bson.M{"Email": email}
	update := bson.M{"$set": bson.M{"Verified": true}}
	_, err = UsersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false
	}

	return true
}

// SetUserRole updates the role of the user with the given user ID.
func SetUserRole(userID primitive.ObjectID, role string) error {
	client, err := MongoConnection()
	UsersCollection := client.Database("UserDatabase").Collection("Users")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"Role": role}}

	_, err = UsersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

// GetUserData retrieves the user data for the given email.
func GetUserData(email string) (model.ApplicationUser, error) {
	client, err := MongoConnection()
	UsersCollection := client.Database("UserDatabase").Collection("Users")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	filter := bson.M{"Email": email}

	var result model.ApplicationUser

	err = UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return model.ApplicationUser{}, err
	}

	return result, nil
}

func SetAuthorise(userID primitive.ObjectID, authorised bool) error {
	client, err := MongoConnection()
	UsersCollection := client.Database("UserDatabase").Collection("Users")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"Authorised": authorised}}

	_, err = UsersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
