package models

import (
    "time"

    "github.com/google/uuid"
)

type User struct {
    ID            uuid.UUID  `json:"id" db:"id"`
    XsollaUserID  string     `json:"xsolla_user_id" db:"xsolla_user_id"`
    Username      string     `json:"username" db:"username"`
    Email         string     `json:"email" db:"email"`
    TokenBalance  int        `json:"token_balance" db:"token_balance"`
    LastLoginAt   *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
    CreatedAt     time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

type UserProfile struct {
    ID           uuid.UUID       `json:"id"`
    Username     string          `json:"username"`
    Email        string          `json:"email"`
    TokenBalance int             `json:"token_balance"`
    LastLoginAt  *time.Time      `json:"last_login_at,omitempty"`
    Inventory    []InventoryItem `json:"inventory"`
}