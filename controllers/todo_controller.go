package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"todo-api-fiber/async"
	"todo-api-fiber/extensions"
	"todo-api-fiber/managers"
	"todo-api-fiber/models"
	. "todo-api-fiber/models"
	"todo-api-fiber/repositories"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var todoRepository = repositories.BaseRepository{Collection: "todos"}
var todoManager = managers.BaseManager{Collection: "todos"}

// Get all todos
func GetTodosAsync(w http.ResponseWriter, r *http.Request) {
	var params QueryParameters
	decoder.Decode(&params, r.URL.Query())

	future := async.Invoke(func() interface{} {
		cursor, error := todoRepository.GetAllAsync(params.SetValues("todos"))

		if error != nil {
			return false
		}

		var todos []Todo

		cursor.All(context.TODO(), &todos)
		return todos
	})

	val := future.Await()

	if val == false {
		failed := FailResponse{}
		failed.SendInternalServerError(w, r)
		return
	}

	countDocsFuture := async.Invoke(func() interface{} {
		totalDocs, _ := todoRepository.CountDocumentsAsync(params.SetFilter())
		return totalDocs
	})

	count := countDocsFuture.Await()
	offset, limit := params.SetPagination()
	pager := GetPager(int(count.(int64)), offset, limit)

	success := SuccessResponseWithPager{Data: val, Pager: pager}
	success.SendOk(w, r)
}

// Get todos by user ID
func GetTodosByUserAsync(w http.ResponseWriter, r *http.Request) {
	var params models.QueryParameters
	decoder.Decode(&params, r.URL.Query())

	claims := extensions.GetClaims(r)
	params.Filter += "user=" + claims["id"].(string)

	future := async.Invoke(func() interface{} {
		cursor, error := todoRepository.GetAllAsync(params.SetValues("todos"))

		if error != nil {
			return false
		}

		var todos []models.Todo

		cursor.All(context.TODO(), &todos)
		return todos
	})

	val := future.Await()

	if val == false {
		failed := FailResponse{}
		failed.SendInternalServerError(w, r)
		return
	}

	countDocsFuture := async.Invoke(func() interface{} {
		totalDocs, _ := todoRepository.CountDocumentsAsync(params.SetFilter())
		return totalDocs
	})

	count := countDocsFuture.Await()
	offset, limit := params.SetPagination()
	pager := models.GetPager(int(count.(int64)), offset, limit)

	success := SuccessResponseWithPager{Data: val, Pager: pager}
	success.SendOk(w, r)
}

// Get todo by its ID
func GetTodoByIdAsync(w http.ResponseWriter, r *http.Request) {
	future := async.Invoke(func() interface{} {
		var todo models.Todo
		id := mux.Vars(r)["id"]

		singleResult := todoRepository.GetByIdAsync(id)
		singleResult.Decode(&todo)

		return todo
	})

	val := future.Await()

	success := SuccessResponse{Data: val}
	success.SendOk(w, r)
}

// Create new todo
func CreateTodoAsync(w http.ResponseWriter, r *http.Request) {
	var todo models.CreateTodoDto
	json.NewDecoder(r.Body).Decode(&todo)

	structError := structValidator.Struct(todo)

	if structError != nil {
		invalidModel := FailResponse{Message: structError.Error()}
		invalidModel.SendUnprocessableEntity(w, r)
		return
	}

	future := async.Invoke(func() interface{} {
		inserted, error := todoManager.CreateAsync(todo.SetDefaultValues())

		if error != nil {
			return false
		}

		return inserted
	})

	val := future.Await()

	if val == false {
		failed := FailResponse{}
		failed.SendInternalServerError(w, r)
		return
	}

	success := SuccessResponse{Data: val}
	success.SendCreated(w, r)
}

// Creates a new todo using user ID
func CreateTodoByUserAsync(w http.ResponseWriter, r *http.Request) {
	claims := extensions.GetClaims(r)
	var todo models.CreateTodoDto

	json.NewDecoder(r.Body).Decode(&todo)
	user, _ := primitive.ObjectIDFromHex(claims["id"].(string))
	todo.User = user

	structError := structValidator.Struct(todo)

	if structError != nil {
		invalidModel := FailResponse{Message: structError.Error()}
		invalidModel.SendUnprocessableEntity(w, r)
		return
	}

	future := async.Invoke(func() interface{} {
		inserted, error := todoManager.CreateAsync(todo.SetDefaultValues())

		if error != nil {
			return false
		}

		return inserted
	})

	val := future.Await()

	if val == false {
		failed := FailResponse{}
		failed.SendInternalServerError(w, r)
		return
	}

	success := SuccessResponse{Data: val}
	success.SendCreated(w, r)
}

// Update an existing todo
func UpdateTodoByIdAsync(w http.ResponseWriter, r *http.Request) {
	var todo models.UpdateTodoDto
	id := mux.Vars(r)["id"]
	json.NewDecoder(r.Body).Decode(&todo)

	structError := structValidator.Struct(todo)

	if structError != nil {
		invalidModel := FailResponse{Message: structError.Error()}
		invalidModel.SendUnprocessableEntity(w, r)
		return
	}

	future := async.Invoke(func() interface{} {
		singleResult := todoManager.UpdateByIdAsync(id, todo.SetUpdatedDate())
		singleResult.Decode(&todo)

		return todo
	})

	val := future.Await()

	success := SuccessResponse{Data: val}
	success.SendOk(w, r)
}

// Delete an existing todo
func DeleteTodoByIdAsync(w http.ResponseWriter, r *http.Request) {
	future := async.Invoke(func() interface{} {
		id := mux.Vars(r)["id"]
		deleted, error := todoManager.DeleteByIdAsync(id)

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
