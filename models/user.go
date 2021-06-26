package models

import (
	"time"
	"todo-api-fiber/enums"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty"`
	Roles     []string           `json:"roles,omitempty" bson:"roles,omitempty"`
	Todos     []Todo             `json:"todos,omitempty" bson:"todos,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type CreateUserDto struct {
	Name      string    `json:"name" bson:"name" validate:"required"`
	Email     string    `json:"email" bson:"email" validate:"required,email"`
	Password  string    `json:"password" bson:"password" validate:"required"`
	Roles     []string  `json:"roles" bson:"roles" validate:"omitempty,validateRoles"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type UpdateUserDto struct {
	Name      string    `json:"name,omitempty" bson:"name,omitempty"`
	Password  string    `json:"password,omitempty" bson:"password,omitempty"`
	Roles     []string  `json:"roles,omitempty" bson:"roles,omitempty" validate:"omitempty,validateRoles"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// Sets the default values for User instance
func (user *CreateUserDto) SetDefaultValues() *CreateUserDto {
	if user.Roles == nil {
		roles := []string{enums.Client}
		user.Roles = roles
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	user.Password = string(hashed)
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()

	return user
}

// Sets the updated date to the document
func (user *UpdateUserDto) SetUpdatedDate() *UpdateUserDto {
	if len(user.Password) != 0 {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.Password = string(hashed)
	}

	user.UpdatedAt = time.Now().UTC()

	return user
}

// Returns the relationships of User model
func GetUserRelationships() map[string]map[string]primitive.D {
	return map[string]map[string]primitive.D{
		"todos": {
			"lookupStage": bson.D{{
				Key: "$lookup", Value: bson.D{
					{Key: "from", Value: "todos"},
					{Key: "localField", Value: "_id"},
					{Key: "foreignField", Value: "user"},
					{Key: "as", Value: "todos"},
				},
			}},
			"unwindStage": nil,
		},
	}
}
