package routes

import (
	"net/http"
	"todo-api-fiber/controllers"
	"todo-api-fiber/enums"
	"todo-api-fiber/middlewares"

	"github.com/gorilla/mux"
)

var onlySuperAdmin = map[interface{}]bool{
	enums.SuperAdmin: true,
}

var all = map[interface{}]bool{
	enums.All: true,
}

var superAdminAuthorizer = middlewares.IsAuthorized(onlySuperAdmin)
var allAuthorizer = middlewares.IsAuthorized(all)

// Returns the routes for user model
func GetUserRoutes(base *mux.Router) {
	user := base.PathPrefix("/users").Subrouter()
	user.HandleFunc("", controllers.GetUsersAsync).Methods(http.MethodGet)
	user.Use(middlewares.IsAuthenticated)
	user.Use(superAdminAuthorizer)

	userPath := base.PathPrefix("/users/{id}").Subrouter()
	userPath.HandleFunc("", controllers.GetUserByIdAsync).Methods(http.MethodGet)
	userPath.HandleFunc("", controllers.UpdateUserByIdAsync).Methods(http.MethodPatch)
	userPath.HandleFunc("", controllers.DeleteUserByIdAsync).Methods(http.MethodDelete)
	userPath.Use(middlewares.IsAuthenticated)
	userPath.Use(allAuthorizer)
	userPath.Use(middlewares.ExistsUserByIdAsync)
}
