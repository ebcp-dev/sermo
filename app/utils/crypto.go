package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// Used for validating header tokens.
var mySigningKey = []byte(viper.GetString("SIGNING_KEY"))

func HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func ValidateToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return []byte(mySigningKey), nil
	})
	if err != nil {
		return false
	}
	return token.Valid
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	// Return false if no errors.
	return err == nil
}

// Generate JWT
func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = "sermoapi"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	if os.Getenv("ENV") == "prod" {
		mySigningKey = []byte(os.Getenv("SIGNING_KEY"))
	}

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		// fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}
