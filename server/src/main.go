package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Jacobsmi/CarTracker/server/src/dbutils"
	"github.com/Jacobsmi/CarTracker/server/src/dbutils/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

type response struct {
	Success bool
	Err     string
}

type userResponse struct {
	Success  bool
	ID       int
	Name     string
	Username string
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Cors for options requests
func corsOptionsResp(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusAccepted)
}

func corsJsonResp(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

// Sign up Route handler
func signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		corsOptionsResp(w)
	} else if r.Method == "POST" {
		// Decode the JSON from the request into a user object
		var newUser models.User

		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			corsJsonResp(w)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response{false, "json_decode_err"})
			return
		}

		// Hash the password using Golang bcrypt package
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if err != nil {
			corsJsonResp(w)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response{false, "hashing_error"})
			return
		}

		// Create a parameterized SQL Statement to insert a new user into the database
		sqlStatement := `INSERT INTO users(name, username, password) VALUES($1, $2, $3);`

		// Attempt to execute the SQL Statement
		_, err = dbutils.DB.Exec(sqlStatement, strings.ToLower(newUser.Name), newUser.Username, string(hashedPassword))

		// Error handling for DB Insert Statement
		// Mostly done to caught a duplicate error
		if err, ok := err.(*pq.Error); ok {
			corsJsonResp(w)
			switch string(err.Code) {
			case "23505":
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(response{false, "duplicate_user"})
			default:
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response{false, "unhandled_db_error"})
			}
			return
		} else if err != nil {
			// Unknown error
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response{false, "unhandled_error"})
			return
		}
		// Control how long the token should be valid for
		expirationTime := time.Now().Add(20 * time.Minute)

		claims := Claims{
			Username: newUser.Username,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
			},
		}
		// Generates the token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Returns the complete signed token
		tokenString, err := token.SignedString(jwtKey)

		if err != nil {
			corsJsonResp(w)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response{false, "token_err"})
			return
		}
		// Set the cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Expires:  expirationTime,
			HttpOnly: true,
			Secure:   true,
		})
		// Set the headers for the response
		corsJsonResp(w)
		w.WriteHeader(http.StatusAccepted)
		// Set the body of the response
		json.NewEncoder(w).Encode(response{true, ""})
	}
}

func getUserInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		corsOptionsResp(w)
	} else if r.Method == "POST" {
		// Need to check the token here
		tokenCookie, err := r.Cookie("token")
		if err != nil {
			fmt.Println(err)
			corsJsonResp(w)
			json.NewEncoder(w).Encode(response{false, "cookie_parse_error"})
			return
		}
		tokenValue := tokenCookie.Value

		// Create a variable of type claims to hold the claims
		tokenClaims := Claims{}
		// Parse out the token and the claims
		// Token just holds basic info like validity claims have actual token information
		token, err := jwt.ParseWithClaims(tokenValue, &tokenClaims, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
		if err != nil {
			fmt.Println("error parsing token")
			fmt.Println(err)
			return
		}
		// Create a temporary instance of the struct for the SQL data to be written into
		var user models.User
		// Parse out cookie info here to get username
		if token.Valid {
			// Create a SQL statement to get the user info from the database
			sqlStatement := `SELECT * FROM users WHERE username = $1`
			// Queries the database for a single row with user data
			row := dbutils.DB.QueryRow(sqlStatement, tokenClaims.Username)
			// Scan the row that was returned to extract inforamtion
			err := row.Scan(&user.ID, &user.Name, &user.Username, &user.Password)
			if err != nil {
				fmt.Println("Error reading SQL Data")
				fmt.Println(err)
				return
			}
		}
		// Set headers for the response
		corsJsonResp(w)
		json.NewEncoder(w).Encode(userResponse{true, user.ID, user.Name, user.Username})
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	var loginUser models.User

	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		fmt.Println("Error parsing JSON")
		fmt.Println(err)
		corsJsonResp(w)
		json.NewEncoder(w).Encode(response{false, "json_parse_error"})
		return
	}

	sqlStatement := `SELECT * FROM users WHERE username = $1`

	row := dbutils.DB.QueryRow(sqlStatement, loginUser.Username)

	var scannedUser models.User

	row.Scan(&scannedUser.ID, &scannedUser.Name, &scannedUser.Username, &scannedUser.Password)

	samePass := bcrypt.CompareHashAndPassword([]byte(scannedUser.Password), []byte(loginUser.Password))
	if samePass != nil {
		corsJsonResp(w)
		json.NewEncoder(w).Encode(response{false, "wrong_pass"})
		return
	} else {
		// Generate a new token
	}
}

func handleFuncs() {
	http.HandleFunc("/getuserinfo", getUserInfo)
	http.HandleFunc("/signup", signUp)
	http.HandleFunc("/login", login)

	fmt.Println("Running the API at http://localhost:5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func main() {
	defer dbutils.DB.Close()
	handleFuncs()
}
