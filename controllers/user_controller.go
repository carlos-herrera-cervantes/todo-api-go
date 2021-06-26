package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"todo-api-fiber/async"
	"todo-api-fiber/enums"
	"todo-api-fiber/extensions"
	"todo-api-fiber/managers"
	. "todo-api-fiber/models"
	"todo-api-fiber/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var userRepository = repositories.BaseRepository{Collection: "users"}
var userManager = managers.BaseManager{Collection: "users"}
var decoder = schema.NewDecoder()
var structValidator = validator.New()

// Sets the ID depending the user role
func setIdByRole(r *http.Request) string {
	claims := extensions.GetClaims(r)
	roles := claims["role"].([]interface{})
	id := ""

	for i := range roles {
		if roles[i].(string) == enums.SuperAdmin {
			id = mux.Vars(r)["id"]
			break
		}
	}

	if len(id) == 0 {
		id = claims["id"].(string)
	}

	return id
}

// Get all users
func GetUsersAsync(w http.ResponseWriter, r *http.Request) {
	var params QueryParameters
	decoder.Decode(&params, r.URL.Query())

	future := async.Invoke(func() interface{} {
		cursor, error := userRepository.GetAllAsync(params.SetValues("users"))

		if error != nil {
			return false
		}

		var users []User
		cursor.All(context.TODO(), &users)

		return users
	})

	val := future.Await()

	if val == false {
		failed := FailResponse{}
		failed.SendInternalServerError(w, r)
		return
	}

	countDocsFuture := async.Invoke(func() interface{} {
		totalDocs, _ := userRepository.CountDocumentsAsync(params.SetFilter())
		return totalDocs
	})

	count := countDocsFuture.Await()
	offset, limit := params.SetPagination()
	pager := GetPager(int(count.(int64)), offset, limit)
	success := SuccessResponseWithPager{Data: val, Pager: pager}

	success.SendOk(w, r)
}

// Get user by its ID
func GetUserByIdAsync(w http.ResponseWriter, r *http.Request) {
	future := async.Invoke(func() interface{} {
		var params QueryParameters

		decoder.Decode(&params, r.URL.Query())
		params.Filter = "_id=" + setIdByRole(r)

		cursor, _ := userRepository.GetAllAsync(params.SetValues("users"))
		var users []User
		cursor.All(context.TODO(), &users)

		if len(users) == 0 {
			return nil
		}

		return users[0]
	})

	val := future.Await()

	success := SuccessResponse{Data: val}
	success.SendOk(w, r)
}

// Create new user
func CreateUserAsync(w http.ResponseWriter, r *http.Request) {
	var user CreateUserDto
	json.NewDecoder(r.Body).Decode(&user)

	structValidator.RegisterValidation("validateRoles", extensions.ValidateRoleAttribute)
	structError := structValidator.Struct(user)

	if structError != nil {
		invalidModel := FailResponse{Message: structError.Error()}
		invalidModel.SendUnprocessableEntity(w, r)
		return
	}

	future := async.Invoke(func() interface{} {
		inserted, error := userManager.CreateAsync(user.SetDefaultValues())

		if error != nil {
			return false
		}

		return inserted
	})

	val := future.Await()

	if val != false {
		success := SuccessResponse{Data: val}
		success.SendCreated(w, r)
		return
	}

	failed := FailResponse{}
	failed.SendInternalServerError(w, r)
}

// Update an existing user
func UpdateUserByIdAsync(w http.ResponseWriter, r *http.Request) {
	var user UpdateUserDto
	json.NewDecoder(r.Body).Decode(&user)

	structValidator.RegisterValidation("validateRoles", extensions.ValidateRoleAttribute)
	structError := structValidator.Struct(user)

	if structError != nil {
		invalidModel := FailResponse{Message: structError.Error()}
		invalidModel.SendUnprocessableEntity(w, r)
		return
	}

	future := async.Invoke(func() interface{} {
		singleResult := userManager.UpdateByIdAsync(setIdByRole(r), user.SetUpdatedDate())
		singleResult.Decode(&user)

		return user
	})

	val := future.Await()

	success := SuccessResponse{Data: val}
	success.SendOk(w, r)
}

// Delete an existing user
func DeleteUserByIdAsync(w http.ResponseWriter, r *http.Request) {
	future := async.Invoke(func() interface{} {
		deleted, error := userManager.DeleteByIdAsync(setIdByRole(r))

		if error != nil {
			return false
		}

		return deleted
	})

	val := future.Await()

	if val == false {
		failed := FailResponse{}
		failed.SendInternalServerError(w, r)
		return
	}

	success := SuccessResponse{}
	success.SendNoContent(w, r)
}
