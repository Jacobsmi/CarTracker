package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Jacobsmi/CarTracker/server/src/dbutils"
	"github.com/Jacobsmi/CarTracker/server/src/dbutils/models"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type response struct {
	Success bool
	Err     string
}

func errorHandler(w http.ResponseWriter, err error) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	var resp response
	// If it is a pq type error handle here
	if err, ok := err.(*pq.Error); ok {
		switch string(err.Code) {
		case "23505":
			resp = response{false, "duplicate_user"}
		default:
			resp = response{false, "unhandled_db_error"}
		}
	} else {
		// Unknown error
		resp = response{false, "unhandled_err"}
	}
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(resp)
}

func corsOptionsResp(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusAccepted)
}

func signUp(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method == "OPTIONS" {
		corsOptionsResp(w)
	} else if r.Method == "POST" {
		// Decode the JSON from the request into a user object
		var newUser models.User

		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response{false, "json_decode_err"})
			return
		}

		// Hash the password using Golang bcrypt package
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response{false, "hashing_error"})
			return
		}

		// Create a parameterized SQL Statement to insert a new user into the database
		sqlStatement := `INSERT INTO users(name, username, password) VALUES($1, $2, $3);`

		// Attempt to execute the SQL Statement
		_, err = dbutils.DB.Exec(sqlStatement, newUser.Name, newUser.Username, string(hashedPassword))

		// Insert error handling
		if err != nil {
			errorHandler(w, err)
			return
		}
		// If no errors return
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		json.NewEncoder(w).Encode(response{true, ""})
	}
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
