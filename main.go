package main

import "github.com/gorilla/mux"

func main() {
	r := mux.NewRouter()                            // Declaring the router returning a new router instance
	r.HandleFunc("/signup", signup).Methods("POST") // Handling routes first the ("Route", function).REST METHOD USING
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/protected", TokenVerifyMiddleware(ProtectedEndpoint)).Methods("GET") // This makes a protected endpont once the JWT has been verified
}
