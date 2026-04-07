package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"soulforge-backend/config"
	"soulforge-backend/models"
)

// GetInventory получение инвентаря пользователя
func GetInventory(c *gin.Context) {
	userID := c.GetString("userID")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
		return
	}

	inventoriesCollection := config.GetCollection("inventories")

	var inventory models.Inventory
	err = inventoriesCollection.FindOne(c.Request.Context(), bson.M{"user_id": objID}).Decode(&inventory)
	if err == mongo.ErrNoDocuments {
		// Если инвентаря нет, создаем новый
		inventory = *models.NewInventory(objID)
		_, err = inventoriesCollection.InsertOne(c.Request.Context(), inventory)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании инвентаря"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"inventory": inventory,
	})
}

// ClaimDailyReward получение ежедневной награды
func ClaimDailyReward(c *gin.Context) {
	userID := c.GetString("userID")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
		return
	}

	var req struct {
		Day int `json:"day" binding:"required,min=1,max=7"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dailyRewardsCollection := config.GetCollection("daily_rewards")

	var dailyReward models.DailyReward
	err = dailyRewardsCollection.FindOne(c.Request.Context(), bson.M{"user_id": objID}).Decode(&dailyReward)
	if err == mongo.ErrNoDocuments {
		// Если записи нет, создаем новую
		dailyReward = *models.NewDailyReward(objID)
		_, err = dailyRewardsCollection.InsertOne(c.Request.Context(), dailyReward)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании записи ежедневных наград"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
		return
	}

	// Получаем награду за день
	success, message, reward := dailyReward.ClaimDayReward(req.Day)
	if !success {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	// Обновляем запись в базе
	_, err = dailyRewardsCollection.UpdateOne(c.Request.Context(),
		bson.M{"user_id": objID},
		bson.M{"$set": bson.M{
			"last_claimed": dailyReward.LastClaimed,
			"updated_at":   dailyReward.UpdatedAt,
		}},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения ежедневных наград"})
		return
	}

	// Добавляем награду в инвентарь
	if reward != nil {
		inventoriesCollection := config.GetCollection("inventories")

		var inventory models.Inventory
		err = inventoriesCollection.FindOne(c.Request.Context(), bson.M{"user_id": objID}).Decode(&inventory)
		if err == nil {
			// Добавляем материалы
			if materials, ok := reward["materials"].(map[models.Rarity]int); ok {
				for rarity, quantity := range materials {
					inventory.AddMaterials(rarity, quantity)
				}
			}

			// Добавляем валюту
			if currency, ok := reward["currency"].(int); ok {
				inventory.AddCurrency(currency)
			}

			// Обновляем инвентарь
			_, err = inventoriesCollection.UpdateOne(c.Request.Context(),
				bson.M{"user_id": objID},
				bson.M{"$set": bson.M{
					"materials":  inventory.Materials,
					"currency":   inventory.Currency,
					"updated_at": inventory.UpdatedAt,
				}},
			)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      message,
		"reward":       reward,
		"daily_reward": dailyReward,
	})
}

// GetDailyRewardConfig получение конфигурации ежедневных наград
func GetDailyRewardConfig(c *gin.Context) {
	config := models.GetDailyRewardConfig()
	c.JSON(http.StatusOK, gin.H{
		"rewards": config,
	})
}

// AddMaterials добавление материалов (для тестирования и админ-функций)
func AddMaterials(c *gin.Context) {
	userID := c.GetString("userID")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
		return
	}

	var req struct {
		Rarity   string `json:"rarity" binding:"required"`
		Quantity int    `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rarity := models.Rarity(req.Rarity)
	if !isValidRarity(rarity) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверная редкость материала"})
		return
	}

	inventoriesCollection := config.GetCollection("inventories")

	var inventory models.Inventory
	err = inventoriesCollection.FindOne(c.Request.Context(), bson.M{"user_id": objID}).Decode(&inventory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения инвентаря"})
		return
	}

	// Добавляем материалы
	if !inventory.AddMaterials(rarity, req.Quantity) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при добавлении материалов"})
		return
	}

	// Обновляем инвентарь в базе
	_, err = inventoriesCollection.UpdateOne(c.Request.Context(),
		bson.M{"user_id": objID},
		bson.M{"$set": bson.M{
			"materials":  inventory.Materials,
			"updated_at": inventory.UpdatedAt,
		}},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения инвентаря"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Материалы добавлены",
		"inventory": inventory,
	})
}

// GetCraftingHistory получение истории крафта
func GetCraftingHistory(c *gin.Context) {
	userID := c.GetString("userID")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
		return
	}

	inventoriesCollection := config.GetCollection("inventories")

	var inventory models.Inventory
	err = inventoriesCollection.FindOne(c.Request.Context(), bson.M{"user_id": objID}).Decode(&inventory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения инвентаря"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"history": inventory.CraftingHistory,
	})
}

// GetOwnedItems получение купленных предметов
func GetOwnedItems(c *gin.Context) {
	userID := c.GetString("userID")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
		return
	}

	inventoriesCollection := config.GetCollection("inventories")

	var inventory models.Inventory
	err = inventoriesCollection.FindOne(c.Request.Context(), bson.M{"user_id": objID}).Decode(&inventory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения инвентаря"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": inventory.OwnedItems,
	})
}

// Вспомогательная функция
func isValidRarity(rarity models.Rarity) bool {
	switch rarity {
	case models.Grey, models.Green, models.Blue, models.Purple, models.Gold:
		return true
	default:
		return false
	}
}
