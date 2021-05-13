package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
	"todo-api-fiber/async"
	"todo-api-fiber/extensions"
	"todo-api-fiber/managers"
	"todo-api-fiber/models"
	. "todo-api-fiber/models"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var tokenManager = managers.BaseManager{Collection: "tokens"}

// Get the authentication token
func LoginAsync(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	json.NewDecoder(r.Body).Decode(&credentials)

	future := async.Invoke(func() interface{} {
		var user models.User

		singleResult := userRepository.GetOneAsync(bson.M{"email": credentials.Email})
		error := singleResult.Decode(&user)

		if error != nil {
			return nil
		}

		return user
	}).Await()

	failed := FailResponse{}
	key := "InvalidCredentials"

	if future == nil {
		failed.SendBadRequest(w, r, key)
		return
	}

	parsed := future.(User)
	validCredentials := bcrypt.CompareHashAndPassword([]byte(parsed.Password), []byte(credentials.Password))

	if validCredentials != nil {
		failed.SendBadRequest(w, r, key)
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    parsed.ID.Hex(),
		"email": credentials.Email,
		"role":  parsed.Roles,
		"exp":   time.Now().UTC().Add(60 * time.Second).Unix(),
	})

	secret := []byte(os.Getenv("SECRET_KEY"))
	token, _ := claims.SignedString(secret)
	docToken := AccessToken{Token: token, User: parsed.ID, Role: parsed.Roles}

	structValidator.RegisterValidation("validateRoles", extensions.ValidateRoleAttribute)
	structError := structValidator.Struct(docToken)

	if structError != nil {
		invalidModel := FailResponse{Message: structError.Error()}
		invalidModel.SendUnprocessableEntity(w, r)
		return
	}

	async.Invoke(func() interface{} {
		tokenManager.DeleteManyAsync(bson.M{"user": parsed.ID})
		return nil
	}).Await()

	async.Invoke(func() interface{} {
		tokenManager.CreateAsync(docToken.SetDefaultValues())
		return nil
	}).Await()

	success := SuccessResponse{Data: token}
	success.SendOk(w, r)
}
