package controllers

import (
	"context"
	"todo-api-fiber/async"
	"todo-api-fiber/db"
	"todo-api-fiber/models"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var todoCollection = db.Connect().Collection("todos")

// Get all todos
func GetTodosAsync(ctx *fiber.Ctx) {
	future := async.Invoke(func() interface{} {
		cursor, _ := todoCollection.Find(context.TODO(), bson.D{})

		var todos []models.Todo

		cursor.All(context.TODO(), &todos)
		return todos
	})

	val := future.Await()

	ctx.JSON(fiber.Map{
		"status": true,
		"data":   val,
	})
}

// Get todo by its ID
func GetTodoByIdAsync(ctx *fiber.Ctx) {
	future := async.Invoke(func() interface{} {
		var todo models.Todo

		id, _ := primitive.ObjectIDFromHex(ctx.Params("id"))
		todoCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&todo)

		return todo
	})

	val := future.Await()

	ctx.JSON(fiber.Map{
		"status": true,
		"data":   val,
	})
}

// Create new todo
func CreateTodoAsync(ctx *fiber.Ctx) {
	future := async.Invoke(func() interface{} {
		var todo models.Todo

		ctx.BodyParser(&todo)
		inserted, _ := todoCollection.InsertOne(context.TODO(), todo)

		return inserted
	})

	val := future.Await()

	ctx.Status(201).JSON(fiber.Map{
		"status": true,
		"data":   val,
	})
}

// Update an existing todo
func UpdateTodoByIdAsync(ctx *fiber.Ctx) {
	future := async.Invoke(func() interface{} {
		id, _ := primitive.ObjectIDFromHex(ctx.Params("id"))
		after := options.After
		opt := options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		}

		var todo models.Todo

		ctx.BodyParser(&todo)
		todoCollection.FindOneAndUpdate(context.TODO(), bson.M{"_id": id}, bson.M{"$set": todo}, &opt).Decode(&todo)

		return todo
	})

	val := future.Await()

	ctx.JSON(fiber.Map{
		"status": true,
		"data":   val,
	})
}

// Delete an existing todo
func DeleteTodoByIdAsync(ctx *fiber.Ctx) {
	future := async.Invoke(func() interface{} {
		id, _ := primitive.ObjectIDFromHex(ctx.Params("id"))
		todoCollection.DeleteOne(context.TODO(), bson.M{"_id": id})

		return true
	})

	_ = future.Await()

	ctx.Status(204).JSON(fiber.Map{
		"status": true,
		"data":   nil,
	})
}
