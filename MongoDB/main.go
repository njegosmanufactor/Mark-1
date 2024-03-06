package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"Username"`
	Password string             `bson:"Password"`
}

func main() {
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

	// Accessing the database and collection
	collection := client.Database("UserDatabase").Collection("Users")

	for {
		fmt.Println("Enter 1 to add a user\nEnter 2 to display users\nEnter 3 to exit")
		var input int
		fmt.Println("Enter your choice: ")
		fmt.Scanf("%d\n", &input)

		switch input {
		case 1:
			var username string
			var password string
			fmt.Println("Enter username: ")
			fmt.Scanf("%s\n", &username)
			fmt.Println("Enter password: ")
			fmt.Scanf("%s\n", &password)
			// Creating user instance
			user := User{
				Username: username,
				Password: password,
			}

			// Adding user to the database
			insertResult, err := collection.InsertOne(context.Background(), user)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Added new user with ID:", insertResult.InsertedID)
		case 2:
			// Finding all users in the database
			cursor, err := collection.Find(context.Background(), bson.D{})
			if err != nil {
				log.Fatal(err)
			}
			defer cursor.Close(context.Background())

			// Iterating through the results
			var users []User
			for cursor.Next(context.Background()) {
				var user User
				if err := cursor.Decode(&user); err != nil {
					log.Fatal(err)
				}
				users = append(users, user)
			}

			// Checking if an error occurred while iterating through the results
			if err := cursor.Err(); err != nil {
				log.Fatal(err)
			}

			// Printing Users
			for _, u := range users {
				fmt.Printf("User ID: %s, Username: %s, Password: %s\n", u.ID, u.Username, u.Password)
			}
		case 3:
			os.Exit(0)
		default:
			fmt.Println("Unknown command. Please enter valid choice")
			continue
		}

		break
	}
}
