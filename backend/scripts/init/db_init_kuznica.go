package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"soulforge-backend/models"
)

func main() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("kuznica")
	log.Println("Connected to kuznica database")

	if err := initializeCollections(db); err != nil {
		log.Fatal("Error initializing collections:", err)
	}

	log.Println("Database initialization completed successfully")
}

func initializeCollections(db *mongo.Database) error {
	ctx := context.Background()

	if err := initCraftingRecipes(db, ctx); err != nil {
		return err
	}

	if err := initIndexes(db, ctx); err != nil {
		return err
	}

	if err := createDefaultAdmin(db, ctx); err != nil {
		log.Printf("Warning: Could not create default admin user: %v", err)
	}

	return nil
}

func initCraftingRecipes(db *mongo.Database, ctx context.Context) error {
	craftingCollection := db.Collection("crafting_recipes")

	_, err := craftingCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	recipes := []interface{}{
		models.CraftingRecipe{
			Name:        "Upgrade: Grey to Green",
			Description: "Upgrade grey stones to green stones",
			UpgradeRecipe: models.UpgradeRecipe{
				InputRarity:        models.Grey,
				OutputRarity:       models.Green,
				GuaranteedQuantity: 10,
				SuccessChance:      100,
			},
			ItemRecipes: []models.ItemRecipe{},
			IsActive:    true,
			Season:      "default",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		models.CraftingRecipe{
			Name:        "Upgrade: Green to Blue",
			Description: "Upgrade green stones to blue stones",
			UpgradeRecipe: models.UpgradeRecipe{
				InputRarity:        models.Green,
				OutputRarity:       models.Blue,
				GuaranteedQuantity: 10,
				SuccessChance:      100,
			},
			ItemRecipes: []models.ItemRecipe{},
			IsActive:    true,
			Season:      "default",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		models.CraftingRecipe{
			Name:        "Upgrade: Blue to Purple",
			Description: "Upgrade blue stones to purple stones",
			UpgradeRecipe: models.UpgradeRecipe{
				InputRarity:        models.Blue,
				OutputRarity:       models.Purple,
				GuaranteedQuantity: 10,
				SuccessChance:      100,
			},
			ItemRecipes: []models.ItemRecipe{},
			IsActive:    true,
			Season:      "default",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		models.CraftingRecipe{
			Name:        "Upgrade: Purple to Gold",
			Description: "Upgrade purple stones to gold stones",
			UpgradeRecipe: models.UpgradeRecipe{
				InputRarity:        models.Purple,
				OutputRarity:       models.Gold,
				GuaranteedQuantity: 10,
				SuccessChance:      100,
			},
			ItemRecipes: []models.ItemRecipe{},
			IsActive:    true,
			Season:      "default",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	result, err := craftingCollection.InsertMany(ctx, recipes)
	if err != nil {
		return err
	}

	log.Printf("Inserted %d crafting recipes", len(result.InsertedIDs))
	return nil
}

func initIndexes(db *mongo.Database, ctx context.Context) error {
	log.Println("Creating indexes...")

	userIndexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"username": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"is_active": 1},
		},
		{
			Keys: bson.M{"created_at": 1},
		},
	}

	_, err := db.Collection("users").Indexes().CreateMany(ctx, userIndexes)
	if err != nil {
		return err
	}

	craftingIndexes := []mongo.IndexModel{
		{
			Keys: bson.M{"upgrade_recipe.input_rarity": 1},
		},
		{
			Keys: bson.M{"is_active": 1},
		},
		{
			Keys: bson.M{"season": 1},
		},
	}

	_, err = db.Collection("crafting_recipes").Indexes().CreateMany(ctx, craftingIndexes)
	if err != nil {
		return err
	}

	inventoryIndexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"user_id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"owned_items.sku": 1},
		},
		{
			Keys: bson.M{"updated_at": 1},
		},
	}

	_, err = db.Collection("inventories").Indexes().CreateMany(ctx, inventoryIndexes)
	if err != nil {
		return err
	}

	dailyRewardIndexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"user_id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"last_claimed": 1},
		},
	}

	_, err = db.Collection("daily_rewards").Indexes().CreateMany(ctx, dailyRewardIndexes)
	if err != nil {
		return err
	}

	log.Println("All indexes created successfully")
	return nil
}

func createDefaultAdmin(db *mongo.Database, ctx context.Context) error {
	usersCollection := db.Collection("users")
	count, err := usersCollection.CountDocuments(ctx, bson.M{"username": "admin"})
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("Admin user already exists")
		return nil
	}

	adminUser, err := models.NewUser("admin", "admin@soulforge.com", "admin123")
	if err != nil {
		return err
	}

	result, err := usersCollection.InsertOne(ctx, adminUser)
	if err != nil {
		return err
	}

	adminUser.ID = result.InsertedID.(primitive.ObjectID)

	inventory := models.NewInventory(adminUser.ID)

	_, err = db.Collection("inventories").InsertOne(ctx, inventory)
	if err != nil {
		return err
	}

	dailyReward := models.NewDailyReward(adminUser.ID)
	_, err = db.Collection("daily_rewards").InsertOne(ctx, dailyReward)
	if err != nil {
		return err
	}

	log.Println("Default admin user created successfully")
	log.Println("Username: admin")
	log.Println("Password: admin123")
	log.Println("Email: admin@soulforge.com")

	return nil
}
