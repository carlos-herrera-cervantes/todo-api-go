package main

import (
	"todo-api-fiber/routes"

	"github.com/gofiber/fiber"
)

func main() {
	app := fiber.New()

	api := app.Group("/api")
	v1 := api.Group("/v1")

	routes.ConfigAuthRoutes(v1)
	routes.ConfigUserRoutes(v1)
	routes.ConfigTodosRoutes(v1)

	app.Listen(":3000")
}
