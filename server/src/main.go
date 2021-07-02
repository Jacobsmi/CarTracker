package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Jacobsmi/CarTracker/server/src/dbutils"
)

func signUp(w http.ResponseWriter, r *http.Request) {
	// Write a new user to the database
}

func handleFuncs() {
	http.HandleFunc("/signup", signUp)

	fmt.Println("Running the API at http://localhost:5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func main() {
	defer dbutils.DB.Close()
	handleFuncs()
}
