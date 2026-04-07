package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rarity string

const (
	Grey   Rarity = "grey"
	Green  Rarity = "green"
	Blue   Rarity = "blue"
	Purple Rarity = "purple"
	Gold   Rarity = "gold"
)

type Inventory struct {
	ID              primitive.ObjectID    `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID    `bson:"user_id" json:"user_id"`
	Materials       map[Rarity]int       `bson:"materials" json:"materials"`
	Currency        int                   `bson:"currency" json:"currency"`
	OwnedItems      []OwnedItem           `bson:"owned_items" json:"owned_items"`
	CraftingHistory []CraftingHistory     `bson:"crafting_history" json:"crafting_history"`
	DailyRewards    DailyRewards          `bson:"daily_rewards" json:"daily_rewards"`
	CreatedAt       time.Time             `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time             `bson:"updated_at" json:"updated_at"`
}

type OwnedItem struct {
	Sku        string    `bson:"sku" json:"sku"`
	Name       string    `bson:"name" json:"name"`
	Type       string    `bson:"type" json:"type"` // skin, weapon, currency, etc.
	ImageURL   string    `bson:"image_url" json:"image_url"`
	ObtainedAt time.Time `bson:"obtained_at" json:"obtained_at"`
	Source     string    `bson:"source" json:"source"` // purchase, craft, reward
	Rarity     string    `bson:"rarity" json:"rarity"`
}

type CraftingHistory struct {
	MaterialsUsed map[Rarity]int `bson:"materials_used" json:"materials_used"`
	Result        CraftResult    `bson:"result" json:"result"`
	Timestamp     time.Time      `bson:"timestamp" json:"timestamp"`
}

type CraftResult struct {
	Success       bool            `bson:"success" json:"success"`
	ResultRarity  Rarity          `bson:"result_rarity" json:"result_rarity"`
	ResultQuantity int             `bson:"result_quantity" json:"result_quantity"`
	Cashback      map[Rarity]int  `bson:"cashback" json:"cashback"`
}

type DailyRewards struct {
	LastClaim     time.Time `bson:"last_claim" json:"last_claim"`
	CurrentStreak int       `bson:"current_streak" json:"current_streak"`
	TotalClaimed  int       `bson:"total_claimed" json:"total_claimed"`
}

// NewInventory создает новый инвентарь для пользователя
func NewInventory(userID primitive.ObjectID) *Inventory {
	now := time.Now()
	return &Inventory{
		UserID: userID,
		Materials: map[Rarity]int{
			Grey:   0,
			Green:  0,
			Blue:   0,
			Purple: 0,
			Gold:   0,
		},
		Currency:        0,
		OwnedItems:      []OwnedItem{},
		CraftingHistory: []CraftingHistory{},
		DailyRewards: DailyRewards{
			CurrentStreak: 0,
			TotalClaimed:  0,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// AddMaterials добавляет материалы в инвентарь
func (i *Inventory) AddMaterials(rarity Rarity, quantity int) bool {
	if quantity <= 0 {
		return false
	}
	
	if i.Materials == nil {
		i.Materials = make(map[Rarity]int)
	}
	
	i.Materials[rarity] += quantity
	i.UpdatedAt = time.Now()
	return true
}

// RemoveMaterials удаляет материалы из инвентаря
func (i *Inventory) RemoveMaterials(rarity Rarity, quantity int) bool {
	if quantity <= 0 || i.Materials[rarity] < quantity {
		return false
	}
	
	i.Materials[rarity] -= quantity
	i.UpdatedAt = time.Now()
	return true
}

// AddCurrency добавляет валюту
func (i *Inventory) AddCurrency(amount int) bool {
	if amount <= 0 {
		return false
	}
	
	i.Currency += amount
	i.UpdatedAt = time.Now()
	return true
}

// DeductCurrency списывает валюту
func (i *Inventory) DeductCurrency(amount int) bool {
	if amount <= 0 || i.Currency < amount {
		return false
	}
	
	i.Currency -= amount
	i.UpdatedAt = time.Now()
	return true
}

// AddItem добавляет предмет в инвентарь
func (i *Inventory) AddItem(item OwnedItem) bool {
	// Проверяем, нет ли уже такого предмета
	for _, ownedItem := range i.OwnedItems {
		if ownedItem.Sku == item.Sku {
			return false
		}
	}
	
	item.ObtainedAt = time.Now()
	i.OwnedItems = append(i.OwnedItems, item)
	i.UpdatedAt = time.Now()
	return true
}

// AddCraftingHistory добавляет запись в историю крафта
func (i *Inventory) AddCraftingHistory(history CraftingHistory) {
	history.Timestamp = time.Now()
	i.CraftingHistory = append(i.CraftingHistory, history)
	
	// Ограничиваем историю последними 100 записями
	if len(i.CraftingHistory) > 100 {
		i.CraftingHistory = i.CraftingHistory[len(i.CraftingHistory)-100:]
	}
	
	i.UpdatedAt = time.Now()
}

// CanClaimDailyReward проверяет, можно ли получить ежедневную награду
func (i *Inventory) CanClaimDailyReward() bool {
	if i.DailyRewards.LastClaim.IsZero() {
		return true
	}
	
	now := time.Now()
	lastClaim := i.DailyRewards.LastClaim
	
	// Проверяем, был ли уже клейм сегодня
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	lastClaimDay := time.Date(lastClaim.Year(), lastClaim.Month(), lastClaim.Day(), 0, 0, 0, 0, lastClaim.Location())
	
	return !today.Equal(lastClaimDay)
}

// ClaimDailyReward получает ежедневную награду
func (i *Inventory) ClaimDailyReward() map[string]interface{} {
	if !i.CanClaimDailyReward() {
		return nil
	}
	
	now := time.Now()
	lastClaim := i.DailyRewards.LastClaim
	
	// Базовые награды в зависимости от стрика
	rewards := make(map[string]interface{})
	
	baseMaterials := map[Rarity]int{
		Grey: 3,
		Green: 1,
	}
	baseCurrency := 50
	
	// Бонусы за стрик
	streakBonus := i.DailyRewards.CurrentStreak / 7
	if streakBonus > 0 {
		baseMaterials[Blue] = streakBonus
		baseCurrency += streakBonus * 100
	}
	
	// Обновляем информацию о наградах
	i.DailyRewards.LastClaim = now
	i.DailyRewards.TotalClaimed++
	
	// Проверяем, был ли вход вчера
	yesterday := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
	lastClaimDay := time.Date(lastClaim.Year(), lastClaim.Month(), lastClaim.Day(), 0, 0, 0, 0, lastClaim.Location())
	
	if lastClaimDay.Equal(yesterday) {
		i.DailyRewards.CurrentStreak++
	} else {
		i.DailyRewards.CurrentStreak = 1
	}
	
	// Добавляем награды в инвентарь
	for rarity, quantity := range baseMaterials {
		i.AddMaterials(rarity, quantity)
		rewards[string(rarity)] = quantity
	}
	
	i.AddCurrency(baseCurrency)
	rewards["currency"] = baseCurrency
	
	rewards["streak"] = i.DailyRewards.CurrentStreak
	
	i.UpdatedAt = time.Now()
	
	return rewards
}
