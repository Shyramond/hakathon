package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"soulforge-backend/config"

	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	// Подключаемся к БД
	config.ConnectDB()
	
	ctx := context.Background()
	
	fmt.Println("=== DUMP DATABASE ===")
	
	// Получаем список всех коллекций
	db := config.DB
	collections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		log.Fatal("Error listing collections:", err)
	}
	
	fmt.Printf("Found %d collections\n\n", len(collections))
	
	// Создаем папку для дампа
	dumpDir := "db_dump_" + time.Now().Format("20060102_150405")
	os.MkdirAll(dumpDir, 0755)
	
	for _, collName := range collections {
		dumpCollection(ctx, collName, dumpDir)
	}
	
	fmt.Printf("\n=== DUMP COMPLETED ===\n")
	fmt.Printf("Files saved to: %s/\n", dumpDir)
}

func dumpCollection(ctx context.Context, collName, dumpDir string) {
	collection := config.GetCollection(collName)
	
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error querying %s: %v", collName, err)
		return
	}
	defer cursor.Close(ctx)
	
	var documents []map[string]interface{}
	if err := cursor.All(ctx, &documents); err != nil {
		log.Printf("Error decoding %s: %v", collName, err)
		return
	}
	
	if len(documents) == 0 {
		fmt.Printf("📁 %s: (empty)\n", collName)
		return
	}
	
	// Выводим структуру первого документа
	fmt.Printf("\n📁 %s (%d docs):\n", collName, len(documents))
	printStructure(documents[0], "  ")
	
	// Сохраняем в JSON
	filename := fmt.Sprintf("%s/%s.json", dumpDir, collName)
	data, _ := json.MarshalIndent(documents, "", "  ")
	os.WriteFile(filename, data, 0644)
	fmt.Printf("  → Saved to %s\n", filename)
}

func printStructure(doc map[string]interface{}, indent string) {
	for key, val := range doc {
		switch v := val.(type) {
		case map[string]interface{}:
			fmt.Printf("%s📂 %s:\n", indent, key)
			printStructure(v, indent+"  ")
		case []interface{}:
			fmt.Printf("%s📋 %s: array[%d]\n", indent, key, len(v))
		default:
			fmt.Printf("%s• %s: %T = %v\n", indent, key, val, val)
		}
	}
}
