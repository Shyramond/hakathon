package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DailyReward struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	LastClaimed *time.Time         `bson:"last_claimed,omitempty" json:"last_claimed,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type DailyRewardConfig struct {
	Day    int    `json:"day"`
	Tier   string `json:"tier"`
	Amount int    `json:"amount"`
	Desc   string `json:"desc"`
}

// NewDailyReward создает новую запись ежедневных наград для пользователя
func NewDailyReward(userID primitive.ObjectID) *DailyReward {
	now := time.Now()

	dr := &DailyReward{
		UserID:      userID,
		LastClaimed: nil,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return dr
}

// GetDailyRewardConfig возвращает конфигурацию наград по дням
func GetDailyRewardConfig() []DailyRewardConfig {
	return []DailyRewardConfig{
		{Day: 1, Tier: "grey", Amount: 3, Desc: "A modest start."},
		{Day: 2, Tier: "grey", Amount: 5, Desc: "Gathering more dust."},
		{Day: 3, Tier: "green", Amount: 1, Desc: "First spark of magic!"},
		{Day: 4, Tier: "grey", Amount: 10, Desc: "A handful of stones."},
		{Day: 5, Tier: "green", Amount: 3, Desc: "Feeling lucky."},
		{Day: 6, Tier: "blue", Amount: 1, Desc: "Rare and shiny."},
		{Day: 7, Tier: "purple", Amount: 1, Desc: "The weekly grand prize!"},
	}
}

// CanClaim проверяет, может ли пользователь получить награду за указанный день
func (dr *DailyReward) CanClaim(day int) bool {
	if dr.LastClaimed == nil {
		return true
	}

	now := time.Now()
	lastClaim := *dr.LastClaimed

	// Проверяем, была ли уже получена награда за этот день
	hoursPassed := now.Sub(lastClaim).Hours()

	// Если прошло меньше 24 часов - нельзя получить
	if hoursPassed < 24 {
		return false
	}

	return true
}

// ClaimDayReward получает награду за указанный день
func (dr *DailyReward) ClaimDayReward(day int) (bool, string, map[string]interface{}) {
	if !dr.CanClaim(day) {
		return false, "Награду за этот день уже получена или прошло меньше 24 часов", nil
	}

	// Получаем конфигурацию наград
	config := GetDailyRewardConfig()
	var rewardConfig *DailyRewardConfig

	for _, cfg := range config {
		if cfg.Day == day {
			rewardConfig = &cfg
			break
		}
	}

	if rewardConfig == nil {
		return false, "Награда для этого дня не найдена", nil
	}

	// Создаем награду
	reward := map[string]interface{}{
		"materials": map[Rarity]int{Rarity(rewardConfig.Tier): rewardConfig.Amount},
	}

	// Обновляем время последней награды
	now := time.Now()
	dr.LastClaimed = &now
	dr.UpdatedAt = now

	return true, "Награда успешно получена!", reward
}
