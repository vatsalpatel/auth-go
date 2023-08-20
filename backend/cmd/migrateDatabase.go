package main

import (
	"log"

	"bitbucket.org/vatsal64/va_pa/pkg/storage"
	"github.com/joho/godotenv"
)

func MigrateDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	db := storage.GetMySQLDatabase()
	db.Exec("DROP TABLE IF EXISTS users")
	db.Exec(`CREATE TABLE users 
		(
			id INT AUTO_INCREMENT PRIMARY KEY, 
			email VARCHAR(255) UNIQUE,
			password VARCHAR(255), 
			name VARCHAR(255), 
			phone VARCHAR(255), 
			isGoogleUser BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);`,
	)
}

func main() {
	MigrateDatabase()
}
