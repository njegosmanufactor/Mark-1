package MongoDB

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OtherUser struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"Username"`
}

type ApplicationUser struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"Username"`
	Password string             `bson:"Password"`
}

// save user into database
func SaveUserOther(username string) {
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
	user := OtherUser{
		Username: username,
	}

	// Adding user to the database
	insertResult, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added new user with ID:", insertResult.InsertedID)
}

func SaveUserApplication(username string, password string) {
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
	user := ApplicationUser{
		Username: username,
		Password: password,
	}

	// Adding user to the database
	insertResult, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added new user with ID:", insertResult.InsertedID)
}

func ValidUser(username string, password string) bool {
	uri := "mongodb+srv://nikolakojic:Bombarder535@userdatabase.y6rrj9g.mongodb.net/?retryWrites=true&w=majority&appName=UserDataBase"
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("UserDatabase").Collection("Users")
	filter := bson.M{"Username": username, "Password": password}
	var result ApplicationUser
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatal(err)
	}
	return true
}

func ValidUsername(username string) bool {
	uri := "mongodb+srv://nikolakojic:Bombarder535@userdatabase.y6rrj9g.mongodb.net/?retryWrites=true&w=majority&appName=UserDataBase"
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("UserDatabase").Collection("Users")
	filter := bson.M{"Username": username}
	var result ApplicationUser
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatal(err)
	}
	return true
}
