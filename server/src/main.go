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
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/rs/cors"
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
	ID int `json:"id"`
	jwt.StandardClaims
}

func generateToken(ID int, w http.ResponseWriter) {
	// Control how long the token should be valid for
	expirationTime := time.Now().Add(20 * time.Minute)
	// Generate a new token
	claims := Claims{
		ID: ID,
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
}

// Get User Info
func getUserInfo(w http.ResponseWriter, r *http.Request) {
	// Need to check the token here
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		fmt.Println(err)

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
		sqlStatement := `SELECT * FROM users WHERE id = $1`
		// Queries the database for a single row with user data
		row := dbutils.DB.QueryRow(sqlStatement, tokenClaims.ID)
		// Scan the row that was returned to extract inforamtion
		err := row.Scan(&user.ID, &user.Name, &user.Username, &user.Password)
		if err != nil {
			fmt.Println("Error reading SQL Data")
			fmt.Println(err)
			return
		}
	}
	// Set headers for the response
	json.NewEncoder(w).Encode(userResponse{true, user.ID, user.Name, user.Username})
}

// Sign up Route handler attempts to decode JSON into a user object, hash the password provided in the JSON,
// and then try to insert the new user with hashed password into the database, then issues a token with the new user's ID
func signUp(w http.ResponseWriter, r *http.Request) {

	// Decode the JSON from the request into a user object
	var newUser models.User

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response{false, "json_decode_err"})
		return
	}

	// Hash the password using Golang bcrypt package
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response{false, "hashing_error"})
		return
	}

	// Create a parameterized SQL Statement to insert a new user into the database
	sqlStatement := `INSERT INTO users(name, username, password) VALUES($1, $2, $3) RETURNING ID;`

	// Attempt to execute the SQL Statement
	_, err = dbutils.DB.Exec(sqlStatement, strings.ToLower(newUser.Name), newUser.Username, string(hashedPassword))
	// Error handling for DB Insert Statement
	// Mostly done to caught a duplicate error
	if err, ok := err.(*pq.Error); ok {
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

	// Get the ID of the newly created user so it can be used as a claim in the new JWT
	var id int
	sqlStatement = `SELECT id FROM users WHERE username = $1`
	// Queries the database for a single row with user data
	row := dbutils.DB.QueryRow(sqlStatement, newUser.Username)
	// Scan the row that was returned to extract inforamtion
	err = row.Scan(&id)
	if err != nil {
		fmt.Println("Error reading SQL Data")
		fmt.Println(err)
		return
	}
	// Create a token
	generateToken(id, w)
	// Set the headers for the respons
	w.WriteHeader(http.StatusAccepted)
	// Set the body of the response
	json.NewEncoder(w).Encode(response{true, ""})

}

// Login route gets a username and password from JSON in a POST request
// It then queries the database for the selected user, reads that info in, and compares the passwords
func login(w http.ResponseWriter, r *http.Request) {
	// Create a user to hold the data from the JSON request
	var loginUser models.User

	// Decode the JSON data into the user object
	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		json.NewEncoder(w).Encode(response{false, "json_parse_error"})
		return
	}

	// Prepare the query for a user
	sqlStatement := `SELECT * FROM users WHERE username = $1`

	// Get the user with correspoding username from the database
	row := dbutils.DB.QueryRow(sqlStatement, loginUser.Username)
	if row == nil {
		json.NewEncoder(w).Encode(response{false, "user_not_exist"})
		return
	}
	// Create a user to hold the data from the queried user
	var scannedUser models.User
	// Read the data from the SQL query into the user object with Scan
	row.Scan(&scannedUser.ID, &scannedUser.Name, &scannedUser.Username, &scannedUser.Password)

	// Compare the password provided with the hashed password to see if matches
	// If error is returned that means wrong pass otherwise issue the user a JWT
	samePass := bcrypt.CompareHashAndPassword([]byte(scannedUser.Password), []byte(loginUser.Password))
	if samePass != nil {
		json.NewEncoder(w).Encode(response{false, "wrong_pass"})
		return
	} else {
		// Create a token
		generateToken(scannedUser.ID, w)
		// Set the headers for the response
		w.WriteHeader(http.StatusAccepted)
		// Set the body of the response
		json.NewEncoder(w).Encode(response{true, ""})

	}
}

func handleFuncs() {
	r := mux.NewRouter()

	r.HandleFunc("/getuserinfo", getUserInfo)
	r.HandleFunc("/signup", signUp)
	r.HandleFunc("/login", login)

	fmt.Println("Running the API at http://localhost:5000")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)
	log.Fatal(http.ListenAndServe(":5000", handler))
}

func main() {
	defer dbutils.DB.Close()
	handleFuncs()
}
