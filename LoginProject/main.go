package main

import (
	controller "LoginProject/Controller"
	"net/http"
)

func main() {
	http.HandleFunc("/register", controller.RegisterHandler)
	http.HandleFunc("/users", controller.GetUsersHandler)
	http.HandleFunc("/login", controller.LoginHandler)
	http.HandleFunc("/changePassword", controller.ChangePasswordHandler)

	http.ListenAndServe(":8080", nil)
}
