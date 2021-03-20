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

var todoCollection = db.Connect().Collection("todos")

// Verifies if the to-do exists by the ID specified
func ExistsTodoById(ctx *fiber.Ctx) {
	future := async.Invoke(func() interface{} {
		var todo collection.Todo

		id, _ := primitive.ObjectIDFromHex(ctx.Params("id"))
		result := todoCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&todo)

		return result
	})

	lang := ctx.Get("Accept-Language")
	message := locales.GetLocalizer(lang, "TodoNotFound")
	val := future.Await()

	if val != nil {
		ctx.Status(404).JSON(fiber.Map{
			"status":  false,
			"code":    "TodoNotFound",
			"message": message,
		})
		return
	}

	ctx.Next()
}
