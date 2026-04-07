package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

// ConnectDB подключается к MongoDB
func ConnectDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Ошибка подключения к MongoDB:", err)
	}

	// Проверка подключения
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Не удалось подключиться к MongoDB:", err)
	}

	DB = client.Database("kuznica")
	log.Println("Connected to MongoDB: kuznica")
}

// GetCollection возвращает коллекцию
func GetCollection(name string) *mongo.Collection {
	return DB.Collection(name)
}
