package routes

import (
	"todo-api-fiber/controllers"
	"todo-api-fiber/middlewares"

	"github.com/gofiber/fiber"
)

// Sets the endpoints for to-do model
func ConfigTodosRoutes(v1 fiber.Router) {
	v1.Use("/todos/:id", middlewares.ExistsTodoById)

	v1.Get("/todos", controllers.GetTodosAsync)
	v1.Get("/todos/:id", controllers.GetTodoByIdAsync)
	v1.Post("/todos", controllers.CreateTodoAsync)
	v1.Patch("/todos/:id", controllers.UpdateTodoByIdAsync)
	v1.Delete("/todos/:id", controllers.DeleteTodoByIdAsync)
}
