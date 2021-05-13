package repositories

import (
	"context"
	"todo-api-fiber/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BaseRepository struct {
	Collection string
}

var database = db.Connect()

// Returns a mongo cursor towards collection
func (base *BaseRepository) GetAllAsync(pipeline mongo.Pipeline) (*mongo.Cursor, error) {
	collection := database.Collection(base.Collection)
	cursor, error := collection.Aggregate(context.TODO(), pipeline)
	return cursor, error
}

// Returns a single result towards a document
func (base *BaseRepository) GetByIdAsync(id string) *mongo.SingleResult {
	parsed, _ := primitive.ObjectIDFromHex(id)
	collection := database.Collection(base.Collection)
	singleResult := collection.FindOne(context.TODO(), bson.M{"_id": parsed})

	return singleResult
}

// Returns a single result towards a document
func (base *BaseRepository) GetOneAsync(filter interface{}) *mongo.SingleResult {
	collection := database.Collection(base.Collection)
	singleResult := collection.FindOne(context.TODO(), filter)
	return singleResult
}

// Returns the number of documents
func (base *BaseRepository) CountDocumentsAsync(filter interface{}) (int64, error) {
	collection := database.Collection(base.Collection)
	return collection.CountDocuments(context.TODO(), filter)
}
