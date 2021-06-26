package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccessToken struct {
	Token     string             `bson:"token" validate:"required"`
	User      primitive.ObjectID `bson:"user" validate:"required"`
	Role      []string           `bson:"role" validate:"required,validateRoles"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}

// Sets the default values for Token instance
func (token *AccessToken) SetDefaultValues() *AccessToken {
	token.CreatedAt = time.Now().UTC()
	token.UpdatedAt = time.Now().UTC()
	return token
}
