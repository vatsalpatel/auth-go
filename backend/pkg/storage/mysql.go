package storage

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
)

var (
	once     sync.Once
	instance *sql.DB
)

// NewDatabase initializes and returns a new Database instance.
func NewDatabase() *sql.DB {
	once.Do(func() {
		// Open a new database connection
		MYSQL_URI := os.Getenv("MYSQL_URI")
		db, err := sql.Open("mysql", MYSQL_URI)
		if err != nil {
			panic(err)
		}

		// Test the database connection
		err = db.Ping()
		if err != nil {
			panic(err)
		}

		fmt.Println("Connected to the database")

		instance = db
	})

	return instance
}
