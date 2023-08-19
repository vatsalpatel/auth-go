package main

import (
	"log"

	"bitbucket.org/vatsal64/va_pa/pkg/storage"
	"github.com/joho/godotenv"
)

func MigrateDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := storage.GetMySQLDatabase()
	db.Exec("DROP TABLE IF EXISTS users")
	db.Exec(`CREATE TABLE users 
		(
			id INT AUTO_INCREMENT PRIMARY KEY, 
			name VARCHAR(255), 
			email VARCHAR(255), 
			password VARCHAR(255), 
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);`,
	)
}

func main() {
	MigrateDatabase()
}
