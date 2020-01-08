package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type user struct { // Model
	ID       int    `json-:"id"`
	Email    string `json-:"email"`
	Password string `json-:"password"`
}

// JWT ...
type JWT struct { // Model
	Token string `json-:"token"`
}

// Error ...
type Error struct { // Model
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
	json.NewDecoder(r.Body).Decode(&user) // Reciving

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
	spew.Dump(user) // Printing out the user object
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
	w.Header().Set("Content-Type", "application/json") // Setting a header for the response
	responseJSON(w, user)
}

// GenerateToken ...
func GenerateToken(user user) (string, error) {
	var err error
	secret := "secret"

	/*
		jwt - header.payload.secret
		(algorithm, struct containing claims)
		iss - is a type of claim
	*/

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss":   "course",
	})

	// Getting the signed token
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}

func signin(w http.ResponseWriter, r *http.Request) {
	var user user   // Holding user information
	var jwt JWT     // token information
	var error Error // Error information

	json.NewDecoder(r.Body).Decode(&user)

	if user.Email == "" {
		error.Message = "Email is empty."
		respondWithError(w, http.StatusBadRequest, error)
		return
	}
	if user.Password == "" {
		error.Message = "Password is empty."
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	// Store password so that we can later compare with what is saved in out table
	password := user.Password
	// Checking if the user exists in our table * QueryRow returns at least one row
	row := db.QueryRow("select * from users where email=$1", user.Email)

	err := row.Scan(&user.ID, &user.Email, &user.Password)

	// Considiration if the scan get's no rows
	if err != nil {
		if err == sql.ErrNoRows {
			error.Message = "The user does not exist."
			respondWithError(w, http.StatusBadRequest, error)
			return
		} else {
			log.Fatal(err)
		}
	}
	hashedPassword := user.Password
	// Comparing password with pasword saved when logged in turn password into bytes to compare
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		error.Message = "Invalid Password"
		respondWithError(w, http.StatusUnauthorized, error)
		return
	}
	// If there is no error we generate the token
	token, err := GenerateToken(user)

	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK) // Add ok status to the header
	jwt.Token = token            // Assign generated token to the token variable we created

	// Send the jwt token to the user using the function we created earlier
	responseJSON(w, jwt)

}

func protectedEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Protected Endpoint invoked.")
}

// * Validates the token that we sent and gives us access to the protected endpoint
// next is the function that will be called after the token has been verified
func tokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// Pass in a callback function
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errorObject Error                       // Created so any error we encounter we can send to the client
		authHeader := r.Header.Get("Authorization") // Map of key value pairs
		// Extract the body out of the bearer token
		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) == 2 {
			authToken := bearerToken[1]
			// Parses and validates and returns a token
			token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				// Validate the token
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return []byte("secret"), nil // Authorized token
			})
			if error != nil {
				errorObject.Message = error.Error()
				respondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}
			// Token is valid
			if token.Valid {
				next.ServeHTTP(w, r)
			} else {
				errorObject.Message = error.Error()
				respondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}

		} else {
			errorObject.Message = "Invalid Token."
			respondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}

	})
}
