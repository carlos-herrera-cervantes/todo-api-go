package routes

import (
	"net/http"
	"todo-api-fiber/controllers"

	"github.com/gorilla/mux"
)

// Returns the routes for authentication process
func GetAuthRoutes(base *mux.Router) {
	login := base.PathPrefix("/auth").Subrouter()
	login.HandleFunc("/login", controllers.LoginAsync).Methods(http.MethodPost)
	login.HandleFunc("/sign-up", controllers.CreateUserAsync).Methods(http.MethodPost)
}
