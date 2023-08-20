package helpers

import (
	"errors"
	"os"
	"time"

	"bitbucket.org/vatsal64/va_pa/internal/models"
	"bitbucket.org/vatsal64/va_pa/pkg/storage"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type MyClaims struct {
	UserID int    `json:"userID"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func ValidateLogin(email, password string) (models.User, error) {
	db := storage.GetMySQLDatabase()

	var user models.User
	// Query the database to retrieve the hashed password for the given email.
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&hashedPassword)
	if err != nil {
		return user, err
	}

	// Compare the stored hashed password with the provided plain text password.
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// Passwords do not match; return an error or handle the authentication failure as needed.
		return user, err
	}

	// Passwords match; retrieve the user's additional information.
	err = db.QueryRow("SELECT * FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Phone, &user.IsGoogleUser, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

func ValidateGoogleLogin(email string) (models.User, error) {
	db := storage.GetMySQLDatabase()

	var user models.User
	db.QueryRow("SELECT * FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Phone, &user.IsGoogleUser, &user.CreatedAt, &user.UpdatedAt)

	return user, nil
}

func RegisterUser(email, password, name, phone string, isGoogleUser bool) (models.User, error) {
	db := storage.GetMySQLDatabase()

	var user models.User
	// Generate a bcrypt hashed password from the user's plain password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}

	res, err := db.Exec("INSERT INTO users (email, password, name, phone, isGoogleUser) VALUES (?, ?, ?, ?, ?);", email, hashedPassword, name, phone, isGoogleUser)
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
		"userID": user.ID,
		"email":  user.Email,
	})
	jwt.NewWithClaims(jwt.SigningMethodHS256, MyClaims{
		UserID: user.ID,
		Email:  user.Email,
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
	return claims, nil
}

func GetUserByID(userId int) models.User {
	db := storage.GetMySQLDatabase()
	var user models.User
	db.QueryRow("SELECT * FROM users WHERE id = ?", userId).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Phone, &user.IsGoogleUser, &user.CreatedAt, &user.UpdatedAt)
	return user
}

func UpdateUserByID(userId int, name, email, phone string) error {
	db := storage.GetMySQLDatabase()

	var userFromDb models.User
	db.QueryRow("SELECT email, isGoogleUser FROM users where id = ?", userId).Scan(&userFromDb.Email, &userFromDb.IsGoogleUser)
	if userFromDb.IsGoogleUser {
		email = userFromDb.Email
	}

	res, err := db.Exec("UPDATE users SET name = ?, email = ?, phone = ? WHERE id = ?;", name, email, phone, userId)
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
