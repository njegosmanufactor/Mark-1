package MongoDB

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"Username"`
	Name     string             `bson:"Name"`
	Surname  string             `bson:"Surname"`
}

// save user into database
func SaveUser(username string, name string, surname string) {
	// Setting up the URL to connect to the MongoDB server
	uri := "mongodb+srv://nikolakojic:Bombarder535@userdatabase.y6rrj9g.mongodb.net/?retryWrites=true&w=majority&appName=UserDataBase"

	// Setting up client options for connection
	clientOptions := options.Client().ApplyURI(uri)

	// Connecting to the MongoDB server
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Connection to MongoDB successful")
	collection := client.Database("UserDatabase").Collection("Users")

	// Creating user instance
	user := User{
		Username: username,
		Name:     name,
		Surname:  surname,
	}

	// Adding user to the database
	insertResult, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added new user with ID:", insertResult.InsertedID)
}
