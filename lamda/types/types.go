package types

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username     string `json:"username" dynamodbav:"username"`
	PasswordHash string `json:"password" dynamodbav:"password"`
}

func NewUser(u RegisterUser) (User, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return User{}, err
	}
	return User{
		Username:     u.Username,
		PasswordHash: string(passHash),
	}, nil
}

func ValidatePassword(hashedPassword, plainTextPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTextPassword))
}

func CreateToken(user User) string {
	now := time.Now()
	validUntil := now.Add(1 * time.Hour).Unix()

	claims := jwt.MapClaims{
		"user":    user.Username,
		"expires": validUntil,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
		nil,
	)
	secret := "this_is_a_super_secure_secret"
	tokeString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Printf("error %v\n", err)
		return ""
	}
	return tokeString
}
