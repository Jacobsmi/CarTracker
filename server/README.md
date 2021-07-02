# Car Tracker API

- Provides backend functionality for the CarTracker frontend
- Must be run locally for the project to work

## Running the Server

- Step 1
  - Create a PostgreSQL database called `cartracker`
- Step 2
  - Create a file called `dbvars.go` in `server/src/dbutils/` with the following content
  ```go
  package dbutils

  const (
	host     = "localhost"
	port     = 5432
	user     = "your=db-user"
	password = "your-pass"
	dbname   = "cartracker"
  )
  ```
  and customize the variable values for your database
- Step 3
  - Change into the `migrations` directory and run the `main.go` in the migrations directory to set up the tables necessary for the server to run
- Step 4
  - Change into the `src` directory and run the `main.go` and the server will be ready for use
