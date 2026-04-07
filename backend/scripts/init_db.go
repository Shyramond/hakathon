package main

import (
	"context"
	"log"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
	"soulforge-backend/models"
)

func main() {
	// Подключение к MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Ошибка подключения к MongoDB:", err)
	}
	defer client.Disconnect(context.Background())
	
	db := client.Database("soulforge")
	log.Println("Подключено к базе данных soulforge")
	
	// Создание базовых рецептов крафта
	craftingCollection := db.Collection("crafting_recipes")
	
	// Рецепт для улучшения серых камней
	greyToGreen := models.CraftingRecipe{
		Name:        "Улучшение: Серый → Зеленый",
		Description: "Улучшение серых камней до зеленых",
		UpgradeRecipe: models.UpgradeRecipe{
			InputRarity:       models.Grey,
			OutputRarity:      models.Green,
			GuaranteedQuantity: 10,
			SuccessChance:     100,
		},
		ItemRecipes: []models.ItemRecipe{
			{
				Sku:      "skin_basic",
				Name:     "Базовый скин",
				Type:     "skin",
				Rarity:   "common",
				ImageURL: "https://example.com/skin_basic.png",
				Materials: map[models.Rarity]int{
					models.Green: 5,
					models.Blue:  2,
				},
				CurrencyCost: 500,
			},
		},
		IsActive: true,
		Season:   "default",
	}
	
	// Рецепт для улучшения зеленых камней
	greenToBlue := models.CraftingRecipe{
		Name:        "Улучшение: Зеленый → Синий",
		Description: "Улучшение зеленых камней до синих",
		UpgradeRecipe: models.UpgradeRecipe{
			InputRarity:       models.Green,
			OutputRarity:      models.Blue,
			GuaranteedQuantity: 10,
			SuccessChance:     100,
		},
		ItemRecipes: []models.ItemRecipe{
			{
				Sku:      "skin_rare",
				Name:     "Редкий скин",
				Type:     "skin",
				Rarity:   "rare",
				ImageURL: "https://example.com/skin_rare.png",
				Materials: map[models.Rarity]int{
					models.Blue:   3,
					models.Purple: 1,
				},
				CurrencyCost: 1500,
			},
		},
		IsActive: true,
		Season:   "default",
	}
	
	// Рецепт для улучшения синих камней
	blueToPurple := models.CraftingRecipe{
		Name:        "Улучшение: Синий → Фиолетовый",
		Description: "Улучшение синих камней до фиолетовых",
		UpgradeRecipe: models.UpgradeRecipe{
			InputRarity:       models.Blue,
			OutputRarity:      models.Purple,
			GuaranteedQuantity: 10,
			SuccessChance:     100,
		},
		ItemRecipes: []models.ItemRecipe{
			{
				Sku:      "skin_epic",
				Name:     "Эпический скин",
				Type:     "skin",
				Rarity:   "epic",
				ImageURL: "https://example.com/skin_epic.png",
				Materials: map[models.Rarity]int{
					models.Purple: 5,
					models.Gold:   2,
				},
				CurrencyCost: 5000,
				RequiresSubscription: true,
			},
		},
		IsActive: true,
		Season:   "default",
	}
	
	// Рецепт для улучшения фиолетовых камней
	purpleToGold := models.CraftingRecipe{
		Name:        "Улучшение: Фиолетовый → Золотой",
		Description: "Улучшение фиолетовых камней до золотых",
		UpgradeRecipe: models.UpgradeRecipe{
			InputRarity:       models.Purple,
			OutputRarity:      models.Gold,
			GuaranteedQuantity: 10,
			SuccessChance:     100,
		},
		ItemRecipes: []models.ItemRecipe{
			{
				Sku:      "skin_legendary",
				Name:     "Легендарный скин",
				Type:     "skin",
				Rarity:   "legendary",
				ImageURL: "https://example.com/skin_legendary.png",
				Materials: map[models.Rarity]int{
					models.Gold: 10,
				},
				CurrencyCost: 10000,
				RequiresSubscription: true,
			},
		},
		IsActive: true,
		Season:   "default",
	}
	
	recipes := []interface{}{greyToGreen, greenToBlue, blueToPurple, purpleToGold}
	
	// Очистка существующих рецептов
	_, err = craftingCollection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Ошибка очистки рецептов: %v", err)
	}
	
	// Вставка новых рецептов
	result, err := craftingCollection.InsertMany(context.Background(), recipes)
	if err != nil {
		log.Printf("Ошибка вставки рецептов: %v", err)
	} else {
		log.Printf("Вставлено %d рецептов", len(result.InsertedIDs))
	}
	
	// Создание индексов для улучшения производительности
	indexes := []mongo.IndexModel{
		{
			Keys: bson.M{"user_id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"upgrade_recipe.input_rarity": 1},
		},
		{
			Keys: bson.M{"item_recipes.sku": 1},
		},
		{
			Keys: bson.M{"is_active": 1},
		},
	}
	
	_, err = craftingCollection.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		log.Printf("Ошибка создания индексов: %v", err)
	} else {
		log.Println("Индексы созданы успешно")
	}
	
	log.Println("Инициализация базы данных завершена")
}
