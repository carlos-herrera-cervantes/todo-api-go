package routes

import (
	"net/http"
	"todo-api-fiber/controllers"
	"todo-api-fiber/middlewares"

	"github.com/gorilla/mux"
)

// Returns the routes for todo model
func GetTodoRoutes(base *mux.Router) {
	todo := base.PathPrefix("/todos").Subrouter()
	todo.HandleFunc("", controllers.GetTodosAsync).Methods(http.MethodGet)
	todo.HandleFunc("", controllers.CreateTodoAsync).Methods(http.MethodPost)
	todo.Use(middlewares.IsAuthenticated)
	todo.Use(superAdminAuthorizer)

	todoMe := base.PathPrefix("/todos/me").Subrouter()
	todoMe.HandleFunc("", controllers.GetTodosByUserAsync).Methods(http.MethodGet)
	todoMe.HandleFunc("", controllers.CreateTodoByUserAsync).Methods(http.MethodPost)
	todoMe.Use(middlewares.IsAuthenticated)
	todoMe.Use(allAuthorizer)

	todoPath := base.PathPrefix("/todos/{id}").Subrouter()
	todoPath.HandleFunc("", controllers.GetTodoByIdAsync).Methods(http.MethodGet)
	todoPath.HandleFunc("", controllers.UpdateTodoByIdAsync).Methods(http.MethodPatch)
	todoPath.HandleFunc("", controllers.DeleteTodoByIdAsync).Methods(http.MethodDelete)
	todoPath.Use(middlewares.IsAuthenticated)
	todoPath.Use(middlewares.ExistsTodoById)
	todoPath.Use(middlewares.CheckUserTodo)
	todoPath.Use(allAuthorizer)
}
