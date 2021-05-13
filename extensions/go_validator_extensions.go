package extensions

import (
	"todo-api-fiber/enums"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Validates if the role attribute contains the only valid roles
func ValidateRoleAttribute(fl validator.FieldLevel) bool {
	roles := fl.Field().Interface().([]string)
	validRoles := []string{enums.SuperAdmin, enums.Client}

	for i := range roles {
		if !ArrayIncludes(validRoles, roles[i]) {
			return false
		}
	}

	return true
}

// Validates if the id is a valid object id
func ValidateObjectId(fl validator.FieldLevel) bool {
	id := fl.Field().String()
	_, error := primitive.ObjectIDFromHex(id)

	if error != nil {
		return false
	}

	return true
}
