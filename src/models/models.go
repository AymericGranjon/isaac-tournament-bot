package models

import (
	"database/sql"
	"errors"
	"os"
	"strings"

	// This is the documented way to use the driver
	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

// Models contains a list of objects representing database tables
type Models struct {
	Races
	Racers
}

// Init opens a database connection based on the credentials in the ".env" file
func Init() (*Models, error) {
	// Read the database configuration from environment variables
	// (it was loaded from the .env file in main.go)
	dbHost := os.Getenv("DB_HOST")
	if len(dbHost) == 0 {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if len(dbPort) == 0 {
		// 3306 is the default port for MariaDB
		dbPort = "3306"
	}
	dbUser := os.Getenv("DB_USER")
	if len(dbUser) == 0 {
		return nil, errors.New("The \"DB_USER\" environment variable is blank. Set it in the \".env\" file.")
	}
	dbPass := os.Getenv("DB_PASS")
	if len(dbPass) == 0 {
		return nil, errors.New("The \"DB_PASS\" environment variable is blank. Set it in the \".env\" file.")
	}
	dbName := os.Getenv("DB_NAME")
	if len(dbPass) == 0 {
		return nil, errors.New("The \"DB_NAME\" environment variable is blank. Set it in the \".env\" file.")
	}

	// Initialize the database
	dsn := dbUser + ":" + dbPass + "@(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true"
	if v, err := sql.Open("mysql", dsn); err != nil {
		return nil, err
	} else {
		db = v
	}

	// Create the model
	return &Models{}, nil
}

// Close exposes the ability to close the underlying database connection
func (*Models) Close() {
	db.Close()
}

/*
	Misc. subroutines
*/

func stringToSlice(a string) []string {
	slice := strings.Split(a, ",")
	if a == "" {
		slice = make([]string, 0)
	}
	return slice
}

func sliceToString(slice []string) string {
	return strings.Join(slice, ",")
}
