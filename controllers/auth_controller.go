package controllers

import (
	"os"
	"todo-api-fiber/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	_ "github.com/joho/godotenv/autoload"
)

// Get the authentication token
func Login(ctx *fiber.Ctx) {
	var credentials models.Credentials

	ctx.BodyParser(&credentials)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": credentials.Email,
	})

	secret := []byte(os.Getenv("SECRET_KEY"))
	token, _ := claims.SignedString(secret)

	ctx.JSON(fiber.Map{
		"status": true,
		"data":   token,
	})
}
