package main

import (
	controller "MongoGoogle/Controller"
	conn "MongoGoogle/Repository"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

}

func main() {
	conn.InitConnection()
	controller.Authentication()
}
