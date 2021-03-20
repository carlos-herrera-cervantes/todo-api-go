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

var userCollection = db.Connect().Collection("users")

// Get all users
func GetUsersAsync(ctx *fiber.Ctx) {
	future := async.Invoke(func() interface{} {
		var users []models.User

		cursor, _ := userCollection.Find(context.TODO(), bson.D{})
		cursor.All(context.TODO(), &users)

		return users
	})

	val := future.Await()

	ctx.JSON(fiber.Map{
		"status": true,
		"data":   val,
	})
}

// Get user by its ID
func GetUserByIdAsync(ctx *fiber.Ctx) {
	future := async.Invoke(func() interface{} {
		var user models.User

		id, _ := primitive.ObjectIDFromHex(ctx.Params("id"))
		userCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)

		return user
	})

	val := future.Await()

	ctx.JSON(fiber.Map{
		"status": true,
		"data":   val,
	})
}

// Create new user
func CreateUserAsync(ctx *fiber.Ctx) {
	future := async.Invoke(func() interface{} {
		var user models.User

		ctx.BodyParser(&user)
		inserted, _ := userCollection.InsertOne(context.TODO(), user)

		return inserted
	})

	val := future.Await()

	ctx.Status(201).JSON(fiber.Map{
		"status": true,
		"data":   val,
	})
}

// Update an existing user
func UpdateUserByIdAsync(ctx *fiber.Ctx) {
	future := async.Invoke(func() interface{} {
		id, _ := primitive.ObjectIDFromHex(ctx.Params("id"))
		after := options.After
		opt := options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		}

		var user models.User

		ctx.BodyParser(&user)
		userCollection.FindOneAndUpdate(context.TODO(), bson.M{"_id": id}, bson.M{"$set": user}, &opt).Decode(&user)

		return user
	})

	val := future.Await()

	ctx.JSON(fiber.Map{
		"status": true,
		"data":   val,
	})
}

// Delete an existing user
func DeleteUserByIdAsync(ctx *fiber.Ctx) {
	future := async.Invoke(func() interface{} {
		id, _ := primitive.ObjectIDFromHex(ctx.Params("id"))
		userCollection.DeleteOne(context.TODO(), bson.M{"_id": id})

		return true
	})

	_ = future.Await()

	ctx.Status(204).JSON(fiber.Map{
		"status": true,
		"data":   nil,
	})
}
