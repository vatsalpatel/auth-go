package helpers

import (
	"os"
	"time"

	"bitbucket.org/vatsal64/va_pa/internal/models"
	"bitbucket.org/vatsal64/va_pa/pkg/storage"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
	exp      int64
	iat      int64
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
		"exp":      time.Now().Add(time.Hour * 24 * 10).Unix(),
		"iat":      time.Now().Unix(),
	})

	signed, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func ExtractClaims(token string) (jwt.Claims, error) {
	key := []byte(os.Getenv("JWT_SECRET"))
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	return parsed.Claims, nil
}
