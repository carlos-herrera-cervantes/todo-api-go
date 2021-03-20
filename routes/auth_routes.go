package routes

import (
	"todo-api-fiber/controllers"
	"todo-api-fiber/middlewares"
	"todo-api-fiber/models"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
)

// Sets the endpoints for authentication process
func ConfigAuthRoutes(v1 fiber.Router) {
	v1.Use("/auth/login", func(ctx *fiber.Ctx) {
		var user models.User

		ctx.BodyParser(&user)

		filter := bson.M{"email": user.Email}
		middlewares.ExistsUserByFilterAsync(ctx, filter)
	})

	v1.Post("/auth/login", controllers.Login)
}
