package models

import (
    "time"

    "github.com/google/uuid"
)

type TransactionType string

const (
    TransactionLoginReward TransactionType = "login_reward"
    TransactionPurchase    TransactionType = "purchase"
    TransactionRefund      TransactionType = "refund"
    TransactionAdmin       TransactionType = "admin"
)

type TokenTransaction struct {
    ID          uuid.UUID       `json:"id" db:"id"`
    UserID      uuid.UUID       `json:"user_id" db:"user_id"`
    Amount      int             `json:"amount" db:"amount"`
    Type        TransactionType `json:"type" db:"type"`
    ReferenceID *uuid.UUID      `json:"reference_id,omitempty" db:"reference_id"`
    Description string          `json:"description" db:"description"`
    CreatedAt   time.Time       `json:"created_at" db:"created_at"`
}