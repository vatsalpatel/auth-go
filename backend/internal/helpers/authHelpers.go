package helpers

import (
	"errors"
	"log"
	"os"
	"time"

	"bitbucket.org/vatsal64/va_pa/internal/models"
	"bitbucket.org/vatsal64/va_pa/pkg/storage"
	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func ValidateLogin(username, password string) (models.User, error) {
	db := storage.GetMySQLDatabase()

	var user models.User
	db.QueryRow("SELECT * FROM users WHERE username = ? AND password = ?", username, password).Scan(&user.ID, &user.Username, &user.Password, &user.Name, &user.Email, &user.Phone, &user.CreatedAt, &user.UpdatedAt)

	return user, nil
}

func RegisterUser(username, password, name, email, phone string) (models.User, error) {
	db := storage.GetMySQLDatabase()

	var user models.User
	res, err := db.Exec("INSERT INTO users (username, password, name, email, phone) VALUES (?, ?, ?, ?, ?);", username, password, name, email, phone)
	if err != nil {
		return user, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return user, err
	}

	user.ID = int(id)
	return user, nil
}

func GenerateToken(user models.User) (string, error) {
	key := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   user.ID,
		"username": user.Username,
	})
	jwt.NewWithClaims(jwt.SigningMethodHS256, MyClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	signed, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func ExtractClaims(token string) (MyClaims, error) {
	key := []byte(os.Getenv("JWT_SECRET"))
	var claims MyClaims
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return MyClaims{}, err
	}
	log.Println(claims)
	return claims, nil
}

func GetUserByID(userId int) models.User {
	db := storage.GetMySQLDatabase()
	var user models.User
	db.QueryRow("SELECT * FROM users WHERE id = ?", userId).Scan(&user.ID, &user.Username, &user.Password, &user.Name, &user.Email, &user.Phone, &user.CreatedAt, &user.UpdatedAt)
	return user
}

func UpdateUserByID(userId int, name, email, phone string) error {
	db := storage.GetMySQLDatabase()
	res, err := db.Exec("UPDATE users SET name = ?, email = ?, phone = ? WHERE id = ?", name, email, phone, userId)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("unable to update user")
	}

	return nil
}
