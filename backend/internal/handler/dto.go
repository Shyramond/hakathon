package handler

import (
    "time"

    "github.com/google/uuid"
    "github.com/Shyramond/hakathon/backend/internal/models"
    "github.com/Shyramond/hakathon/backend/internal/service"
)

type LoginRequest struct {
    Token string `json:"token" binding:"required"`
}

type EquipRequest struct {
    Equip bool `json:"equip"` // true = equip, false = unequip
}

type PaginationQuery struct {
    Limit  int `form:"limit,default=20"`
    Offset int `form:"offset,default=0"`
}

func (p *PaginationQuery) Validate() {
    if p.Limit <= 0 {
        p.Limit = 20
    }
    if p.Limit > 100 {
        p.Limit = 100
    }
    if p.Offset < 0 {
        p.Offset = 0
    }
}

type LoginResponse struct {
    User        UserDTO `json:"user"`
    IsNewUser   bool    `json:"is_new_user"`
    RewardGiven bool    `json:"reward_given"`
    Reward      int     `json:"reward"`
    NewBalance  int     `json:"new_balance"`
    Message     string  `json:"message"`
}

type UserDTO struct {
    ID           uuid.UUID  `json:"id"`
    Username     string     `json:"username"`
    Email        string     `json:"email"`
    TokenBalance int        `json:"token_balance"`
    LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
    CreatedAt    time.Time  `json:"created_at"`
}

type ProfileResponse struct {
    User      UserDTO           `json:"user"`
    Inventory []InventoryItemDTO `json:"inventory"`
}

type BenefitDTO struct {
    ID          uuid.UUID          `json:"id"`
    Name        string             `json:"name"`
    Description string             `json:"description"`
    Type        models.BenefitType `json:"type"`
    PriceTokens int                `json:"price_tokens"`
    ImageURL    string             `json:"image_url"`
    IsActive    bool               `json:"is_active"`
    CreatedAt   time.Time          `json:"created_at"`
}

type InventoryItemDTO struct {
    ID          uuid.UUID          `json:"id"`
    BenefitID   uuid.UUID          `json:"benefit_id"`
    BenefitName string             `json:"benefit_name"`
    BenefitType models.BenefitType `json:"benefit_type"`
    ImageURL    string             `json:"image_url"`
    PurchasedAt time.Time          `json:"purchased_at"`
    IsEquipped  bool               `json:"is_equipped"`
}

type PurchaseResponse struct {
    Item       InventoryItemDTO `json:"item"`
    Benefit    BenefitDTO       `json:"benefit"`
    TokensSpent int             `json:"tokens_spent"`
    NewBalance  int             `json:"new_balance"`
}

type TransactionDTO struct {
    ID          uuid.UUID              `json:"id"`
    Amount      int                    `json:"amount"`
    Type        models.TransactionType `json:"type"`
    ReferenceID *uuid.UUID             `json:"reference_id,omitempty"`
    Description string                 `json:"description"`
    CreatedAt   time.Time              `json:"created_at"`
}

func toUserDTO(u *models.User) UserDTO {
    return UserDTO{
        ID:           u.ID,
        Username:     u.Username,
        Email:        u.Email,
        TokenBalance: u.TokenBalance,
        LastLoginAt:  u.LastLoginAt,
        CreatedAt:    u.CreatedAt,
    }
}

func toLoginResponse(r *service.LoginResult) LoginResponse {
    return LoginResponse{
        User:        toUserDTO(r.User),
        IsNewUser:   r.IsNewUser,
        RewardGiven: r.RewardGiven,
        Reward:      r.Reward,
        NewBalance:  r.NewBalance,
        Message:     r.Message,
    }
}

func toBenefitDTO(b *models.Benefit) BenefitDTO {
    return BenefitDTO{
        ID:          b.ID,
        Name:        b.Name,
        Description: b.Description,
        Type:        b.Type,
        PriceTokens: b.PriceTokens,
        ImageURL:    b.ImageURL,
        IsActive:    b.IsActive,
        CreatedAt:   b.CreatedAt,
    }
}

func toBenefitDTOs(benefits []models.Benefit) []BenefitDTO {
    dtos := make([]BenefitDTO, len(benefits))
    for i, b := range benefits {
        dtos[i] = toBenefitDTO(&b)
    }
    return dtos
}

func toInventoryItemDTO(item *models.InventoryItem) InventoryItemDTO {
    return InventoryItemDTO{
        ID:          item.ID,
        BenefitID:   item.BenefitID,
        BenefitName: item.BenefitName,
        BenefitType: item.BenefitType,
        ImageURL:    item.ImageURL,
        PurchasedAt: item.PurchasedAt,
        IsEquipped:  item.IsEquipped,
    }
}

func toInventoryItemDTOs(items []models.InventoryItem) []InventoryItemDTO {
    dtos := make([]InventoryItemDTO, len(items))
    for i, item := range items {
        dtos[i] = toInventoryItemDTO(&item)
    }
    return dtos
}

func toTransactionDTO(t *models.TokenTransaction) TransactionDTO {
    return TransactionDTO{
        ID:          t.ID,
        Amount:      t.Amount,
        Type:        t.Type,
        ReferenceID: t.ReferenceID,
        Description: t.Description,
        CreatedAt:   t.CreatedAt,
    }
}

func toTransactionDTOs(txs []models.TokenTransaction) []TransactionDTO {
    dtos := make([]TransactionDTO, len(txs))
    for i, t := range txs {
        dtos[i] = toTransactionDTO(&t)
    }
    return dtos
}

func toPurchaseResponse(r *service.PurchaseResult) PurchaseResponse {
    return PurchaseResponse{
        Item:        toInventoryItemDTO(r.InventoryItem),
        Benefit:     toBenefitDTO(r.Benefit),
        TokensSpent: r.TokensSpent,
        NewBalance:  r.NewBalance,
    }
}