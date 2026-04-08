package models

import (
    "time"

    "github.com/google/uuid"
)

type InventoryItem struct {
    ID          uuid.UUID `json:"id" db:"id"`
    UserID      uuid.UUID `json:"user_id" db:"user_id"`
    BenefitID   uuid.UUID `json:"benefit_id" db:"benefit_id"`
    PurchasedAt time.Time `json:"purchased_at" db:"purchased_at"`
    IsEquipped  bool      `json:"is_equipped" db:"is_equipped"`

    BenefitName string      `json:"benefit_name,omitempty" db:"benefit_name"`
    BenefitType BenefitType `json:"benefit_type,omitempty" db:"benefit_type"`
    ImageURL    string      `json:"image_url,omitempty" db:"image_url"`
}