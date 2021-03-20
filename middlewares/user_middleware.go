package middlewares

import (
	"context"
	"todo-api-fiber/async"
	db "todo-api-fiber/db"
	"todo-api-fiber/locales"
	collection "todo-api-fiber/models"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userCollection = db.Connect().Collection("users")

// Verifies if the user exists by the ID specified
func ExistsUserByIdAsync(ctx *fiber.Ctx) {
	future := async.Invoke(func() interface{} {
		var user collection.User

		id, _ := primitive.ObjectIDFromHex(ctx.Params("id"))
		result := userCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)

		return result
	})

	lang := ctx.Get("Accept-Language")
	message := locales.GetLocalizer(lang, "UserNotFound")
	val := future.Await()

	if val != nil {
		ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"code":    "UserNotFound",
			"message": message,
		})
		return
	}

	ctx.Next()
}

// Verifies if the user exists by specific filter criteria
func ExistsUserByFilterAsync(ctx *fiber.Ctx, filter bson.M) {
	future := async.Invoke(func() interface{} {
		var user collection.User
		result := userCollection.FindOne(context.TODO(), filter).Decode(&user)

		return result
	})

	lang := ctx.Get("Accept-Language")
	message := locales.GetLocalizer(lang, "UserNotFound")
	val := future.Await()

	if val != nil {
		ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"code":    "UserNotFound",
			"message": message,
		})
		return
	}

	ctx.Next()
}
