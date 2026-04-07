package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"soulforge-backend/config"
	"soulforge-backend/models"
)

// UpgradeMaterials улучшение материалов
func UpgradeMaterials(c *gin.Context) {
	userID := c.GetString("userID")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
		return
	}

	var req struct {
		Rarity   string `json:"rarity" binding:"required"`
		Quantity int    `json:"quantity" binding:"required,min=1,max=10"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rarity := models.Rarity(req.Rarity)
	if !isValidRarityCraft(rarity) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверная редкость материала"})
		return
	}

	// Получаем инвентарь
	inventoriesCollection := config.GetCollection("inventories")
	var inventory models.Inventory
	err = inventoriesCollection.FindOne(c.Request.Context(), bson.M{"user_id": objID}).Decode(&inventory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения инвентаря"})
		return
	}

	// Получаем рецепты
	recipesCollection := config.GetCollection("crafting_recipes")
	var recipe models.CraftingRecipe
	err = recipesCollection.FindOne(c.Request.Context(), bson.M{
		"upgrade_recipe.input_rarity": rarity,
		"is_active":                   true,
	}).Decode(&recipe)

	if err == mongo.ErrNoDocuments {
		// Создаем базовый рецепт, если его нет
		recipe = *models.NewDefaultCraftingRecipe()
		recipe.UpgradeRecipe.InputRarity = rarity
		recipe.UpgradeRecipe.OutputRarity = getNextRarity(rarity)

		if recipe.UpgradeRecipe.OutputRarity == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Максимальный уровень редкости достигнут"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения рецептов"})
		return
	}

	// Проверяем возможность улучшения
	if !recipe.CanUpgrade(rarity, &inventory, req.Quantity) {
		// Создаем детальный результат ошибки
		errorResult := gin.H{
			"success": false,
			"error":   "Недостаточно материалов для улучшения",
			"required": map[string]interface{}{
				"materials": map[string]int{string(rarity): req.Quantity},
				"available": inventory.Materials,
			},
			"message": fmt.Sprintf("Требуется %d %s, доступно %d", req.Quantity, rarity, inventory.Materials[rarity]),
		}
		c.JSON(http.StatusBadRequest, errorResult)
		return
	}

	// Выполняем улучшение
	result := recipe.PerformUpgrade(rarity, &inventory, req.Quantity)

	// Обновляем инвентарь в базе
	update := bson.M{
		"materials":        inventory.Materials,
		"crafting_history": inventory.CraftingHistory,
		"updated_at":       inventory.UpdatedAt,
	}

	_, err = inventoriesCollection.UpdateOne(c.Request.Context(),
		bson.M{"user_id": objID},
		bson.M{"$set": update},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения инвентаря"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Улучшение выполнено",
		"result":    result,
		"inventory": inventory,
	})
}

// GetCraftingRecipes получение доступных рецептов улучшения
func GetCraftingRecipes(c *gin.Context) {
	// Получаем рецепты
	recipesCollection := config.GetCollection("crafting_recipes")
	cursor, err := recipesCollection.Find(c.Request.Context(), bson.M{"is_active": true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения рецептов"})
		return
	}
	defer cursor.Close(c.Request.Context())

	var recipes []models.CraftingRecipe
	if err = cursor.All(c.Request.Context(), &recipes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения рецептов"})
		return
	}

	// Если рецептов нет, создаем базовые
	if len(recipes) == 0 {
		baseRecipe := models.NewDefaultCraftingRecipe()
		recipes = append(recipes, *baseRecipe)
	}

	c.JSON(http.StatusOK, gin.H{
		"recipes": recipes,
	})
}

// Helper functions
func isValidRarityCraft(rarity models.Rarity) bool {
	switch rarity {
	case models.Grey, models.Green, models.Blue, models.Purple, models.Gold:
		return true
	default:
		return false
	}
}

func getNextRarity(rarity models.Rarity) models.Rarity {
	switch rarity {
	case models.Grey:
		return models.Green
	case models.Green:
		return models.Blue
	case models.Blue:
		return models.Purple
	case models.Purple:
		return models.Gold
	default:
		return ""
	}
}
