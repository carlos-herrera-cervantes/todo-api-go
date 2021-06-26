package extensions

import (
	"net/http"
	"strings"
	"todo-api-fiber/models"

	"github.com/dgrijalva/jwt-go"
)

// Extracts the string token sended in the request header
func GetStringTokenFromAuthorizationHeader(r *http.Request) string {
	authorization := r.Header.Get("Authorization")
	properties := strings.Split(authorization, " ")
	extracted := properties[len(properties)-1]

	return extracted
}

// Returns the user claims
func GetClaims(r *http.Request) jwt.MapClaims {
	token := GetStringTokenFromAuthorizationHeader(r)
	decoded, _ := models.DecodePayload(token)
	claims := decoded.Claims.(jwt.MapClaims)

	return claims
}
