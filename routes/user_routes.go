package routes

import (
	"todo-api-fiber/controllers"
	"todo-api-fiber/middlewares"

	"github.com/gofiber/fiber"
)

// Sets the endpoints for user model
func ConfigUserRoutes(v1 fiber.Router) {
	v1.Use("/users/:id", middlewares.ExistsUserByIdAsync)

	v1.Get("/users", controllers.GetUsersAsync)
	v1.Get("/users/:id", controllers.GetUserByIdAsync)
	v1.Post("/users", controllers.CreateUserAsync)
	v1.Patch("/users/:id", controllers.UpdateUserByIdAsync)
	v1.Delete("/users/:id", controllers.DeleteUserByIdAsync)
}
