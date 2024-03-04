package main

import (
	auth "RegisterAndLogin/controller"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type LoginDto struct {
	Username string
	Password string
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/login", auth.LogIn).Methods("GET")
	router.HandleFunc("/api/register", auth.Register).Methods("POST")

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", router)
}
