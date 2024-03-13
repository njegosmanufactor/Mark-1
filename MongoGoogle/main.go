package main

import (
	controller "MongoGoogle/Controller"
	conn "MongoGoogle/Repository"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
	conn.Uri, _ = os.LookupEnv("MONGO_URI")
	conn.ClientOptions = options.Client().ApplyURI(conn.Uri)
	conn.Client, conn.Err = mongo.Connect(context.Background(), conn.ClientOptions)
}

func main() {
	fmt.Println(conn.Uri)
	controller.Authentication()
}
