package models

import (
	"errors"
	"os"

	"github.com/dgrijalva/jwt-go"
)

// Extracts the payload used to sign the token
func DecodePayload(extracted string) (*jwt.Token, error) {
	decoded, err := jwt.Parse(extracted, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid json web token")
		}

		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	return decoded, err
}
