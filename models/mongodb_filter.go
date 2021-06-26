package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoDBFilter struct {
	Skip     bson.D
	Limit    bson.D
	Sort     bson.D
	Filter   bson.D
	Relation bson.D
	Unwind   bson.D
}

// Returns the relationships depending on the model
func GetRelationModel(model string) map[string]map[string]primitive.D {
	switch model {
	case "users":
		return GetUserRelationships()
	default:
		return GetTodoRelationships()
	}
}
