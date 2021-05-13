package middlewares

import (
	"encoding/json"
	"net/http"
	"todo-api-fiber/async"
	"todo-api-fiber/locales"
	collection "todo-api-fiber/models"
	"todo-api-fiber/repositories"

	"github.com/gorilla/mux"
)

var userRepository = repositories.BaseRepository{Collection: "users"}

// Verifies if the user exists by the ID specified
func ExistsUserByIdAsync(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		future := async.Invoke(func() interface{} {
			var user collection.User
			id := mux.Vars(r)["id"]

			singleResult := userRepository.GetByIdAsync(id)
			error := singleResult.Decode(&user)

			if error != nil {
				return nil
			}

			return user
		})

		val := future.Await()

		if val != nil {
			next.ServeHTTP(w, r)
			return
		}

		lang := r.Header.Get("Accept-Language")
		message := locales.GetLocalizer(lang, "UserNotFound")

		w.WriteHeader(http.StatusNotFound)

		errorResponse := map[string]interface{}{
			"status":  false,
			"code":    "UserNotFound",
			"message": message,
		}
		json.NewEncoder(w).Encode(errorResponse)
		return
	})
}
