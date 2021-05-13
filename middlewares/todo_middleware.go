package middlewares

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"todo-api-fiber/async"
	"todo-api-fiber/enums"
	"todo-api-fiber/extensions"
	"todo-api-fiber/locales"
	collection "todo-api-fiber/models"
	"todo-api-fiber/repositories"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var todoRepository = repositories.BaseRepository{Collection: "todos"}

// Sets the filter depending the user role
func setFilterByRole(r *http.Request) bson.M {
	claims := extensions.GetClaims(r)
	sub := mux.Vars(r)["id"]

	var filter bson.M
	id, _ := primitive.ObjectIDFromHex(sub)

	if claims["role"] == enums.SuperAdmin {
		filter = bson.M{"_id": id}
	} else {
		user, _ := primitive.ObjectIDFromHex(claims["id"].(string))
		filter = bson.M{"_id": id, "user": user}
	}

	return filter
}

// Verifies if the to-do exists by the ID specified
func ExistsTodoById(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		future := async.Invoke(func() interface{} {
			var todo collection.Todo

			singleResult := todoRepository.GetOneAsync(setFilterByRole(r))
			error := singleResult.Decode(&todo)

			if error != nil {
				return nil
			}

			return todo
		})

		val := future.Await()

		if val != nil {
			next.ServeHTTP(w, r)
			return
		}

		lang := r.Header.Get("Accept-Language")
		message := locales.GetLocalizer(lang, "TodoNotFound")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		errorResponse := map[string]interface{}{
			"status":  false,
			"code":    "TodoNotFound",
			"message": message,
		}
		json.NewEncoder(w).Encode(errorResponse)
	})
}

// Checks if the user in the body request is equal to the user in session
func CheckUserTodo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var parsedBody collection.Todo

		json.NewDecoder(r.Body).Decode(&parsedBody)
		bodyBytes, _ := json.Marshal(parsedBody)

		buf, _ := ioutil.ReadAll(bytes.NewReader(bodyBytes))
		emptyObjectId, _ := primitive.ObjectIDFromHex("000000000000000000000000")

		reader := ioutil.NopCloser(bytes.NewBuffer(buf))
		r.Body = reader

		if parsedBody.User == emptyObjectId {
			next.ServeHTTP(w, r)
			return
		}

		claims := extensions.GetClaims(r)
		roles := claims["role"].([]interface{})

		if extensions.ArrayIncludesT(roles, enums.SuperAdmin) {
			next.ServeHTTP(w, r)
			return
		}

		id, _ := primitive.ObjectIDFromHex(claims["id"].(string))

		if id == parsedBody.User {
			next.ServeHTTP(w, r)
			return
		}

		lang := r.Header.Get("Accept-Language")
		message := locales.GetLocalizer(lang, "InvalidOperation")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)

		errorResponse := map[string]interface{}{
			"status":  false,
			"code":    "InvalidOperation",
			"message": message,
		}
		json.NewEncoder(w).Encode(errorResponse)
	})
}
