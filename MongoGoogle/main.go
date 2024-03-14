package main

import (
	controller "MongoGoogle/Controller"
	conn "MongoGoogle/Repository"
	"context"
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
	controller.Mark1()
	defer func() {
		if err := conn.GetClient().Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
}
