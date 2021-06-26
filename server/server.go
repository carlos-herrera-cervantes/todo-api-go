package server

import (
	"log"
	"net/http"
	"todo-api-fiber/middlewares"
	"todo-api-fiber/routes"

	"github.com/gorilla/mux"
)

// Start web server
func Serve() {
	r := mux.NewRouter()
	base := r.PathPrefix("/api/v1").Subrouter()
	base.Use(middlewares.SetSecurityHeaders)

	routes.GetAuthRoutes(base)
	routes.GetUserRoutes(base)
	routes.GetTodoRoutes(base)

	log.Fatal(http.ListenAndServe(":3000", r))
}
