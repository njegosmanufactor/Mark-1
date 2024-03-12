package main

import (
	authentication "MongoGoogle/LoginRegister"
	"context"
	"log"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var uri = "mongodb+srv://Nikola045:Bombarder535@userdatabase.qcrmscd.mongodb.net/?retryWrites=true&w=majority&appName=UserDataBase"

// Setting up client options for connection
var clientOptions = options.Client().ApplyURI(uri)

// Connecting to the MongoDB server
var Client, Err = mongo.Connect(context.Background(), clientOptions)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {
	authentication.Authentication()
}
