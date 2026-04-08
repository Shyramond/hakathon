package models

import (
    "encoding/json"
    "time"

    "github.com/google/uuid"
)

type BenefitType string

const (
    BenefitTypePromotion  BenefitType = "promotion"
    BenefitTypeSkin       BenefitType = "skin"
    BenefitTypeMascotSkin BenefitType = "mascot_skin"
)

type Benefit struct {
    ID          uuid.UUID       `json:"id" db:"id"`
    Name        string          `json:"name" db:"name"`
    Description string          `json:"description" db:"description"`
    Type        BenefitType     `json:"type" db:"type"`
    PriceTokens int             `json:"price_tokens" db:"price_tokens"`
    ImageURL    string          `json:"image_url" db:"image_url"`
    IsActive    bool            `json:"is_active" db:"is_active"`
    Metadata    json.RawMessage `json:"metadata" db:"metadata"`
    CreatedAt   time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

type BenefitFilter struct {
    Type     *BenefitType `form:"type"`
    IsActive *bool        `form:"is_active"`
    MinPrice *int         `form:"min_price"`
    MaxPrice *int         `form:"max_price"`
}