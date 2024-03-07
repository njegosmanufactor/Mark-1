package main

import (
	authentication "MongoGoogle/LoginRegister"
	"log"

	"github.com/joho/godotenv"
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {
	authentication.Authentication()
}
