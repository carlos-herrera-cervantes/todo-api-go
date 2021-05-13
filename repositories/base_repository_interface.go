package repositories

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type IBaseRepository interface {
	getAllAsync(pipeline mongo.Pipeline) *mongo.Cursor
	getByIdAsync(id string) *mongo.SingleResult
	getOneAsync(filter interface{}) *mongo.SingleResult
	countDocumentsAsync(filter interface{}) (int64, error)
}
