package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title,omitempty" bson:"title,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Done        *bool              `json:"done,omitempty" bson:"done,omitempty"`
	User        primitive.ObjectID `json:"user,omitempty" bson:"user,omitempty"`
	UserBind    *User              `json:"user_bind,omitempty" bson:"user_bind,omitempty"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type CreateTodoDto struct {
	Title       string             `json:"title" bson:"title" validate:"required,min=10,max=30"`
	Description string             `json:"description" bson:"description" validate:"required,min=15"`
	Done        *bool              `json:"done" bson:"done"`
	User        primitive.ObjectID `json:"user" bson:"user" validate:"required"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type UpdateTodoDto struct {
	Title       string             `json:"title,omitempty" bson:"title,omitempty" validate:"omitempty,min=10,max=30"`
	Description string             `json:"description,omitempty" bson:"description,omitempty" validate:"omitempty,min=15"`
	Done        *bool              `json:"done,omitempty" bson:"done,omitempty"`
	User        primitive.ObjectID `json:"user,omitempty" bson:"user,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// Sets the default values for Todo instance
func (todo *CreateTodoDto) SetDefaultValues() *CreateTodoDto {
	if todo.Done == nil {
		myFalse := false
		todo.Done = &myFalse
	}

	todo.CreatedAt = time.Now().UTC()
	todo.UpdatedAt = time.Now().UTC()

	return todo
}

// Sets the updated date to the document
func (todo *UpdateTodoDto) SetUpdatedDate() *UpdateTodoDto {
	todo.UpdatedAt = time.Now().UTC()
	return todo
}

// Returns the relationships of Todo model
func GetTodoRelationships() map[string]map[string]primitive.D {
	return map[string]map[string]primitive.D{
		"users": {
			"lookupStage": bson.D{{
				Key: "$lookup", Value: bson.D{
					{Key: "from", Value: "users"},
					{Key: "localField", Value: "user"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "user_bind"},
				},
			}},
			"unwindStage": bson.D{{
				Key: "$unwind", Value: bson.D{
					{Key: "path", Value: "$user_bind"},
				},
			}},
		},
	}
}
