package Repository

import (
	model "MongoGoogle/Model"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Setting up the URL to connect to the MongoDB server
var Uri = "mongodb+srv://Nikola045:Bombarder535@userdatabase.qcrmscd.mongodb.net/?retryWrites=true&w=majority&appName=UserDataBase"
var ClientOptions = options.Client().ApplyURI(uri)
var Client, Err = mongo.Connect(context.Background(), ClientOptions)

// save user into database
func SaveUserOther(email string) {
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

func SaveUserApplication(email string, firstName string, lastName string, phone string, date string, username string, password string, company string, country string, city string, address string) {
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
		Company:     company,
		Country:     country,
		City:        city,
		Address:     address,
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

func SetUserRoleOwner(email string) error {
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
	update := bson.M{"$set": bson.M{"Role": "Owner"}}

	_, err = UsersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
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
