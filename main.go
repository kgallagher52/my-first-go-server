package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID       int    `json-:"id"`
	Email    string `json-:"email"`
	Password string `json-:"password"`
}

type jwt struct {
	Token string `json-:"token"`
}

// Error ...
type Error struct {
	Message string `json-:"message"`
}

// DB variable
var db *sql.DB

func main() {
	// Establishing a connection to the database
	pgURL, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))

	if err != nil {
		log.Fatal(err)
	}

	// Opening connection ("Driver", URL)
	db, err = sql.Open("postgres", pgURL)

	if err != nil {
		log.Fatal(err)
	}

	// Check if the connection has been established or not if ping comes back empty it has been established.

	err = db.Ping()

	// Declaring the router returning a new router instance
	r := mux.NewRouter()

	// Handling routes first the ("Route", function).Methods("REST") METHOD USING
	r.HandleFunc("/signup", signup).Methods("POST")
	r.HandleFunc("/signin", signin).Methods("POST")

	// This makes a protected endpoint once the JWT has been verified
	r.HandleFunc("/protected", tokenVerifyMiddleware(protectedEndpoint)).Methods("GET")

	// Starting the server
	log.Println("Listen on port 8000....")

	/* log.Fatal() - Will stop the program and log the error
	Method that takes two parameters 1. address 2. handler function */
	log.Fatal(http.ListenAndServe(":8000", r))
}

func respondWithError(w http.ResponseWriter, status int, error Error) {
	// Respond with an error
	// Send status bad request
	w.WriteHeader(status) // 400
	json.NewEncoder(w).Encode(error)
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func signup(w http.ResponseWriter, r *http.Request) {
	var user user
	var error Error
	// Decoder that reads body maps it to user struct
	json.NewDecoder(r.Body).Decode(&user)

	if user.Email == "" { // Validation
		error.Message = "Email is missing."
		respondWithError(w, http.StatusBadRequest, error)
		return

	}

	if user.Password == "" { // Validation
		error.Message = "Password is missing."
		respondWithError(w, http.StatusBadRequest, error)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		log.Fatal(err)
	}
	// Server is expecting a string not a slice of bytes which the hash makes
	user.Password = string(hash)
	spew.Dump(user)
	// Creating query string
	stmt := "insert into users (email, password) values($1,$2) RETURNING id;"

	// Insert records into the Database .Scan(&user.ID) - because we are expecting a return of the id
	err = db.QueryRow(stmt, user.Email, user.Password).Scan(&user.ID)
	// Scan suppose to return an error

	if err != nil {
		error.Message = "Server Error"
		respondWithError(w, http.StatusInternalServerError, error)
		return
	}
	// set as an empty string because we are returning the user object to the user
	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	responseJSON(w, user)
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
