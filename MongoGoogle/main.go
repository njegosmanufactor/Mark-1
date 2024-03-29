package main

import (
	controllers "MongoGoogle/Controller"
	conn "MongoGoogle/Repository"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

}

func main() {
	conn.InitConnection()
	r := mux.NewRouter()

	authController := controllers.NewAuthenticationController()
	authController.RegisterRoutes()

	userController := controllers.NewUserController()
	userController.RegisterRoutes()
	companyController := controllers.NewCompanyController()
	companyController.RegisterRoutes()

	r.PathPrefix("/auth").Handler(authController.Router)
	r.PathPrefix("/users").Handler(userController.Router)
	r.PathPrefix("/company").Handler(companyController.Router)

	fmt.Println("[ UP ON PORT 3000 ]")
	err := http.ListenAndServe(":3000", r)
	log.Fatal(err)
	defer func() {
		if err := conn.GetClient().Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
}
