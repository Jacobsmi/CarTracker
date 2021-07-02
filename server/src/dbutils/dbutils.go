package dbutils

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error

	// Attempt to instantiate the database
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error instantiating the Database")
		panic(err)
	}

	// Ping the database and make a connection
	err = DB.Ping()
	if err != nil {
		fmt.Println("Eror connecting to the database")
		panic(err)
	}
}
