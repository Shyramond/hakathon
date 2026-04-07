package models

import (
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CraftingRecipe struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	Description   string             `bson:"description" json:"description"`
	UpgradeRecipe UpgradeRecipe      `bson:"upgrade_recipe" json:"upgrade_recipe"`
	ItemRecipes   []ItemRecipe       `bson:"item_recipes" json:"item_recipes"`
	IsActive      bool               `bson:"is_active" json:"is_active"`
	Season        string             `bson:"season" json:"season"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type UpgradeRecipe struct {
	InputRarity        Rarity `bson:"input_rarity" json:"input_rarity"`
	OutputRarity       Rarity `bson:"output_rarity" json:"output_rarity"`
	GuaranteedQuantity int    `bson:"guaranteed_quantity" json:"guaranteed_quantity"`
	SuccessChance      int    `bson:"success_chance" json:"success_chance"`
}

type ItemRecipe struct {
	Sku          string         `bson:"sku" json:"sku"`
	Name         string         `bson:"name" json:"name"`
	Type         string         `bson:"type" json:"type"`
	Rarity       string         `bson:"rarity" json:"rarity"`
	ImageURL     string         `bson:"image_url" json:"image_url"`
	Materials    map[Rarity]int `bson:"materials" json:"materials"`
	CurrencyCost int            `bson:"currency_cost" json:"currency_cost"`
}

type CraftAttempt struct {
	InputRarity Rarity `json:"input_rarity"`
	Quantity    int    `json:"quantity"`
}

type CraftAttemptResult struct {
	Success        bool           `json:"success"`
	ResultRarity   Rarity         `json:"result_rarity"`
	ResultQuantity int            `json:"result_quantity"`
	Cashback       map[Rarity]int `json:"cashback"`
	Message        string         `json:"message"`
}

// NewDefaultCraftingRecipe создает рецепт по умолчанию для улучшения
func NewDefaultCraftingRecipe() *CraftingRecipe {
	now := time.Now()

	return &CraftingRecipe{
		Name:        "Базовые рецепты улучшения",
		Description: "Рецепты для улучшения материалов",
		UpgradeRecipe: UpgradeRecipe{
			GuaranteedQuantity: 10,
		},
		ItemRecipes: []ItemRecipe{},
		IsActive:    true,
		Season:      "default",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// GetUpgradeChance возвращает шанс успеха для улучшения
func (cr *CraftingRecipe) GetUpgradeChance(quantity int) int {
	if quantity <= 0 {
		return 0
	}

	// Базовая формула: 10% за каждый материал
	chance := quantity * 10

	// Ограничиваем максимум 100%
	if chance > 100 {
		chance = 100
	}

	return chance
}

// CanUpgrade проверяет, можно ли улучшить материалы
func (cr *CraftingRecipe) CanUpgrade(rarity Rarity, inventory *Inventory, quantity int) bool {
	// Проверяем наличие рецепта для этой редкости
	if cr.UpgradeRecipe.InputRarity != rarity {
		return false
	}

	// Проверяем наличие материалов
	if inventory.Materials[rarity] < quantity {
		return false
	}

	// Проверяем, что есть куда улучшать
	if cr.UpgradeRecipe.OutputRarity == "" {
		return false
	}

	return true
}

// PerformUpgrade выполняет улучшение материалов
func (cr *CraftingRecipe) PerformUpgrade(rarity Rarity, inventory *Inventory, quantity int) CraftAttemptResult {
	if !cr.CanUpgrade(rarity, inventory, quantity) {
		return CraftAttemptResult{
			Success: false,
			Message: "Невозможно выполнить улучшение",
		}
	}

	// Удаляем материалы
	inventory.RemoveMaterials(rarity, quantity)

	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Рассчитываем шанс успеха
	chance := cr.GetUpgradeChance(quantity)

	// Логирование для отладки
	fmt.Printf("Крафт: %s x%d, шанс: %d%%\n", rarity, quantity, chance)

	// Генерируем случайное число
	// Используем нормальный генератор случайных чисел
	success := false
	if chance >= 100 {
		success = true
		fmt.Printf("Гарантированный успех (шанс 100%%)\n")
	} else {
		// Генерируем случайное число от 0 до 99
		roll := rand.Intn(100)
		success = roll < chance
		fmt.Printf("Бросок: %d, успех: %t\n", roll, success)
	}

	result := CraftAttemptResult{
		Success:        success,
		ResultRarity:   cr.UpgradeRecipe.OutputRarity,
		ResultQuantity: 1,
		Cashback:       make(map[Rarity]int),
	}

	if success {
		// Успешное улучшение
		inventory.AddMaterials(cr.UpgradeRecipe.OutputRarity, 1)
		result.Message = "Успешное улучшение!"
	} else {
		// Провал - даем кэшбэк
		result.Message = "Улучшение не удалось"

		// Для серых камней НЕ даем кэшбэк никогда
		if rarity == Grey {
			fmt.Printf("Серые камни - кэшбэк не даем\n")
			// Ничего не даем - абуз prevention
		} else {
			fmt.Printf("Другая редкость - нормальный кэшбэк\n")
			// Для остальных уровней нормальный кэшбэк
			cashbackRarity := getPreviousRarity(rarity)
			if cashbackRarity != "" {
				cashbackAmount := quantity / 2
				if cashbackAmount < 1 {
					cashbackAmount = 1
				}

				inventory.AddMaterials(cashbackRarity, cashbackAmount)
				result.Cashback[cashbackRarity] = cashbackAmount
				result.Message += ". Получен кэшбэк: " + string(cashbackRarity) + " x" + fmt.Sprintf("%d", cashbackAmount)
			}
		}
	}

	// Добавляем запись в историю
	history := CraftingHistory{
		MaterialsUsed: map[Rarity]int{rarity: quantity},
		Result: CraftResult{
			Success:        result.Success,
			ResultRarity:   result.ResultRarity,
			ResultQuantity: result.ResultQuantity,
			Cashback:       result.Cashback,
		},
	}
	inventory.AddCraftingHistory(history)

	return result
}

// CanCraftItem проверяет, можно ли скрафтить предмет
func (cr *CraftingRecipe) CanCraftItem(itemRecipe ItemRecipe, inventory *Inventory) bool {
	for rarity, quantity := range itemRecipe.Materials {
		if inventory.Materials[rarity] < quantity {
			return false
		}
	}

	if inventory.Currency < itemRecipe.CurrencyCost {
		return false
	}

	return true
}

// CraftItem крафтит предмет
func (cr *CraftingRecipe) CraftItem(itemRecipe ItemRecipe, inventory *Inventory) (bool, string) {
	if !cr.CanCraftItem(itemRecipe, inventory) {
		return false, "Недостаточно ресурсов"
	}

	for rarity, quantity := range itemRecipe.Materials {
		inventory.RemoveMaterials(rarity, quantity)
	}

	if itemRecipe.CurrencyCost > 0 {
		inventory.DeductCurrency(itemRecipe.CurrencyCost)
	}

	ownedItem := OwnedItem{
		Sku:      itemRecipe.Sku,
		Name:     itemRecipe.Name,
		Type:     itemRecipe.Type,
		ImageURL: itemRecipe.ImageURL,
		Source:   "craft",
		Rarity:   itemRecipe.Rarity,
	}

	if inventory.AddItem(ownedItem) {
		return true, "Предмет успешно скрафчен"
	}

	return false, "Ошибка при добавлении предмета"
}

// Вспомогательные функции
func getPreviousRarity(rarity Rarity) Rarity {
	switch rarity {
	case Green:
		return Grey
	case Blue:
		return Green
	case Purple:
		return Blue
	case Gold:
		return Purple
	default:
		return ""
	}
}

func getNextRarity(rarity Rarity) Rarity {
	switch rarity {
	case Grey:
		return Green
	case Green:
		return Blue
	case Blue:
		return Purple
	case Purple:
		return Gold
	default:
		return ""
	}
}
