package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"soulforge-backend/config"
	"soulforge-backend/middleware"
	"soulforge-backend/models"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Register регистрация нового пользователя
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, не существует ли пользователь
	usersCollection := config.GetCollection("users")

	var existingUser models.User
	err := usersCollection.FindOne(c.Request.Context(), bson.M{
		"$or": []bson.M{
			{"username": req.Username},
			{"email": req.Email},
		},
	}).Decode(&existingUser)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким именем или email уже существует"})
		return
	}

	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
		return
	}

	// Создаем нового пользователя
	newUser, err := models.NewUser(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании пользователя"})
		return
	}

	// Сохраняем в базу данных
	result, err := usersCollection.InsertOne(c.Request.Context(), newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении пользователя"})
		return
	}

	// Создаем инвентарь для пользователя
	userID := result.InsertedID.(primitive.ObjectID)
	inventory := models.NewInventory(userID)

	inventoriesCollection := config.GetCollection("inventories")
	_, err = inventoriesCollection.InsertOne(c.Request.Context(), inventory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании инвентаря"})
		return
	}

	// Генерируем JWT токен
	token, err := middleware.GenerateJWT(userID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации токена"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Пользователь успешно зарегистрирован",
		"user":    newUser.ToJSON(),
		"token":   token,
	})
}

// Login вход пользователя
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usersCollection := config.GetCollection("users")

	var user models.User
	err := usersCollection.FindOne(c.Request.Context(), bson.M{
		"$or": []bson.M{
			{"username": req.Username},
			{"email": req.Username},
		},
	}).Decode(&user)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверное имя пользователя или пароль"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
		return
	}

	// Проверяем пароль
	if err := user.ComparePassword(req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверное имя пользователя или пароль"})
		return
	}

	// Обновляем информацию о входе
	if user.UpdateDailyLogin() {
		// Если это новый день, обновляем в базе
		usersCollection.UpdateOne(c.Request.Context(),
			bson.M{"_id": user.ID},
			bson.M{"$set": bson.M{
				"last_login":   user.LastLogin,
				"login_streak": user.LoginStreak,
				"total_logins": user.TotalLogins,
				"updated_at":   user.UpdatedAt,
			}},
		)
	}

	// Генерируем JWT токен
	token, err := middleware.GenerateJWT(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Вход выполнен успешно",
		"user":    user.ToJSON(),
		"token":   token,
	})
}

// GetProfile получение профиля пользователя
func GetProfile(c *gin.Context) {
	userID := c.GetString("userID")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
		return
	}

	usersCollection := config.GetCollection("users")

	var user models.User
	err = usersCollection.FindOne(c.Request.Context(), bson.M{"_id": objID}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user.ToJSON(),
	})
}
