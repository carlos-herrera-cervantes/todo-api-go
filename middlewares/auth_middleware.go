package middlewares

import (
	"net/http"
	"todo-api-fiber/async"
	"todo-api-fiber/enums"
	"todo-api-fiber/extensions"
	"todo-api-fiber/models"
	. "todo-api-fiber/models"
	"todo-api-fiber/repositories"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
)

var tokenRepository = repositories.BaseRepository{Collection: "tokens"}

// Checks if the request is authenticated by a user
func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extensions.GetStringTokenFromAuthorizationHeader(r)

		future := async.Invoke(func() interface{} {
			finded := tokenRepository.GetOneAsync(bson.M{"token": token})
			var decoded AccessToken
			err := finded.Decode(&decoded)

			if err != nil {
				return false
			}

			return true
		}).Await()

		failed := FailResponse{}

		if future == false {
			failed.SendForbidden(w, r, "InvalidToken")
			return
		}

		_, err := DecodePayload(token)

		if err == nil {
			next.ServeHTTP(w, r)
			return
		}

		failed.SendForbidden(w, r, "MissingToken")
		return
	})
}

// Checks if the user is authorized to execute actions on the resource
func IsAuthorized(roles map[interface{}]bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			decoded, _ := models.DecodePayload(extensions.GetStringTokenFromAuthorizationHeader(r))

			if roles[enums.All] {
				next.ServeHTTP(w, r)
				return
			}

			claims := decoded.Claims.(jwt.MapClaims)
			userRoles := claims["role"].([]interface{})

			for i := range userRoles {
				finded := roles[userRoles[i]]

				if finded {
					next.ServeHTTP(w, r)
					return
				}
			}

			failed := FailResponse{}
			failed.SendForbidden(w, r, "InvalidPermissions")
			return
		})
	}
}
