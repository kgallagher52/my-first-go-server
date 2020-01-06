package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type user struct {
	ID       int    `json-:"id"`
	Email    string `json-:"email"`
	Password string `json-:"password"`
}

type jwt struct {
	Token string `json-:"token"`
}

type error struct {
	Message string `json-:"message"`
}

func main() {
	// Declaring the router returning a new router instance
	r := mux.NewRouter()
	// Handling routes first the ("Route", function).REST METHOD USING
	r.HandleFunc("/signup", signup).Methods("POST")
	r.HandleFunc("/signin", signin).Methods("POST")
	// This makes a protected endpont once the JWT has been verified
	r.HandleFunc("/protected", tokenVerifyMiddleware(protectedEndpoint)).Methods("GET")

	// log.Fatal() - Will stop the program and log the error
	// Starting the server
	log.Println("Listen on port 8000....")
	log.Fatal(http.ListenAndServe(":8000", r)) // Method that takes two parameters 1. address 2. handler function
}

func signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Signup invoked.")
	w.Write([]byte("Successfully called signup")) // Sending a response
}

func signin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Signin invoked.")
	w.Write([]byte("Successfully called signin")) // Sending a response

}

func protectedEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Protected Endpoint invoked.")
}

func tokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	fmt.Println("token Verify Middleware Invoked.")
	return nil
}
