package main

import (
	"log"
	"os"
	
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	
	"soulforge-backend/config"
	"soulforge-backend/routes"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используются переменные окружения")
	}
	
	// Подключение к MongoDB
	config.ConnectDB()
	
	// Настройка режима Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	// Настройка маршрутов
	router := routes.SetupRoutes()
	
	// Получение порта
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	// Запуск сервера
	log.Printf("Сервер запускается на порту %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
