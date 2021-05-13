package managers

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IBaseManager interface {
	createAsync(document interface{}) (*mongo.InsertOneResult, error)
	updateByIdAsync(id string, document interface{}) *mongo.SingleResult
	deleteByIdAsync(id string) (*mongo.DeleteResult, error)
	deleteManyAsync(filter bson.M) (*mongo.DeleteResult, error)
}
