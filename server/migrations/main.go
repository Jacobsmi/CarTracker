package main

import (
	"fmt"

	"github.com/Jacobsmi/CarTracker/server/src/dbutils"
)

func main() {
	defer dbutils.DB.Close()

	sqlStatement := `CREATE TABLE IF NOT EXISTS users(
		id SERIAL,
		name VARCHAR NOT NULL,
		username VARCHAR UNIQUE NOT NULL,
		password VARCHAR NOT NULL
	)`

	_, err := dbutils.DB.Exec(sqlStatement)
	if err != nil {
		fmt.Println("Error creating table")
		panic(err)
	}

	fmt.Println("Migrations completed successfully")
}
