package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
	Done        bool               `json:"done" bson:"done,omitempty" default0:"false"`
}
