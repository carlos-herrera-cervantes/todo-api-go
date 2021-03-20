package db

import (
	"context"
	"log"
	"os"
	"sync"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientInstance *mongo.Client
var mongoOnce sync.Once

// Connect to MongoDB Server
func Connect() *mongo.Database {
	uri := os.Getenv("MONGODB_HOST") + os.Getenv("MONGODB_DATABASE")

	mongoOnce.Do(func() {
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

		if err != nil {
			log.Fatal(err)
		}

		clientInstance = client
	})

	log.Println("Successfully conected to MongoDB")

	return clientInstance.Database("todo-api-fiber")
}
