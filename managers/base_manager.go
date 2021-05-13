package managers

import (
	"context"
	"todo-api-fiber/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseManager struct {
	Collection string
}

var database = db.Connect()

// Inserts a new document into the collection
func (base *BaseManager) CreateAsync(document interface{}) (*mongo.InsertOneResult, error) {
	collection := database.Collection(base.Collection)
	inserted, error := collection.InsertOne(context.TODO(), document)
	return inserted, error
}

// Updates an existing document into the collection
func (base *BaseManager) UpdateByIdAsync(id string, document interface{}) *mongo.SingleResult {
	collection := database.Collection(base.Collection)
	parsed, _ := primitive.ObjectIDFromHex(id)
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	updated := collection.FindOneAndUpdate(context.TODO(), bson.M{"_id": parsed}, bson.M{"$set": document}, &opt)

	return updated
}

// Deletes an existing document into the collection
func (base *BaseManager) DeleteByIdAsync(id string) (*mongo.DeleteResult, error) {
	collection := database.Collection(base.Collection)
	parsed, _ := primitive.ObjectIDFromHex(id)

	deleted, error := collection.DeleteOne(context.TODO(), bson.M{"_id": parsed})

	return deleted, error
}

// Deletes a set of documents
func (base *BaseManager) DeleteManyAsync(filter bson.M) (*mongo.DeleteResult, error) {
	collection := database.Collection(base.Collection)
	deleted, error := collection.DeleteMany(context.TODO(), filter)
	return deleted, error
}
