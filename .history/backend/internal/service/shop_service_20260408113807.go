package service

import (
    "context"
    "errors"
    "fmt"
    "log/slog"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/Shyramond/hakathon/backend/internal/models"
    "github.com/Shyramond/hakathon/backend/internal/repository"
)

type ShopService struct {
    benefitRepo   repository.BenefitRepository
    inventoryRepo repository.InventoryRepository
    tokenService  *TokenService
    txManager     repository.TxManager
}

type PurchaseResult struct {
    InventoryItem *models.InventoryItem `json:"inventory_item"`
    Benefit       *models.Benefit       `json:"benefit"`
    TokensSpent   int                   `json:"tokens_spent"`
    NewBalance    int                   `json:"new_balance"`
}

func NewShopService(
    benefitRepo repository.BenefitRepository,
    inventoryRepo repository.InventoryRepository,
    tokenService *TokenService,
    txManager repository.TxManager,
) *ShopService {
    return &ShopService{
        benefitRepo:   benefitRepo,
        inventoryRepo: inventoryRepo,
        tokenService:  tokenService,
        txManager:     txManager,
    }
}

// Purchase покупает бенефит для пользователя.
// Вся операция выполняется в одной транзакции:
// 1. Проверить что бенефит существует и активен
// 2. Проверить что пользователь ещё не владеет этим бенефитом
// 3. Проверить и списать баланс
// 4. Добавить в инвентарь
// 5. Записать транзакцию
func (s *ShopService) Purchase(ctx context.Context, userID, benefitID uuid.UUID) (*PurchaseResult, error) {
    benefit, err := s.benefitRepo.GetByID(ctx, benefitID)
    if err != nil {
        if errors.Is(err, repository.ErrBenefitNotFound) {
            return nil, ErrBenefitNotFound
        }
        return nil, fmt.Errorf("get benefit: %w", err)
    }

    if !benefit.IsActive {
        return nil, ErrBenefitNotActive
    }

    var result PurchaseResult
    result.Benefit = benefit
    result.TokensSpent = benefit.PriceTokens

    err = s.txManager.RunInTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
        owned, err := s.inventoryRepo.HasBenefitTx(ctx, tx, userID, benefitID)
        if err != nil {
            return fmt.Errorf("check ownership: %w", err)
        }
        if owned {
            return ErrAlreadyOwned
        }

        description := fmt.Sprintf("Purchase: %s", benefit.Name)
        newBalance, err := s.tokenService.DeductTokens(
            ctx, tx, userID, benefit.PriceTokens,
            models.TransactionPurchase, &benefitID, description,
        )
        if err != nil {
            return err
        }
        result.NewBalance = newBalance

        item, err := s.inventoryRepo.AddTx(ctx, tx, userID, benefitID)
        if err != nil {
            if errors.Is(err, repository.ErrAlreadyOwned) {
                return ErrAlreadyOwned
            }
            return fmt.Errorf("add to inventory: %w", err)
        }
        result.InventoryItem = item

        return nil
    })

    if err != nil {
        return nil, err
    }

    slog.Info("purchase completed",
        "user_id", userID,
        "benefit_id", benefitID,
        "benefit_name", benefit.Name,
        "tokens_spent", benefit.PriceTokens,
        "new_balance", result.NewBalance,
    )

    return &result, nil
}

func (s *ShopService) GetInventory(ctx context.Context, userID uuid.UUID) ([]models.InventoryItem, error) {
    items, err := s.inventoryRepo.GetByUserID(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("get inventory: %w", err)
    }
    if items == nil {
        items = []models.InventoryItem{}
    }
    return items, nil
}

// EquipItem экипирует предмет инвентаря.
func (s *ShopService) EquipItem(ctx context.Context, userID, itemID uuid.UUID) (*models.InventoryItem, error) {
    item, err := s.inventoryRepo.GetByID(ctx, itemID)
    if err != nil {
        if errors.Is(err, repository.ErrInventoryItemNotFound) {
            return nil, ErrItemNotFound
        }
        return nil, fmt.Errorf("get item: %w", err)
    }

    if item.UserID != userID {
        return nil, ErrNotItemOwner
    }

    if item.IsEquipped {
        return nil, ErrAlreadyEquipped
    }

    if err := s.inventoryRepo.UnequipAllByType(ctx, userID, item.BenefitType); err != nil {
        return nil, fmt.Errorf("unequip others: %w", err)
    }

    if err := s.inventoryRepo.SetEquipped(ctx, itemID, true); err != nil {
        return nil, fmt.Errorf("equip item: %w", err)
    }

    item.IsEquipped = true
    return item, nil
}

func (s *ShopService) UnequipItem(ctx context.Context, userID, itemID uuid.UUID) (*models.InventoryItem, error) {
    item, err := s.inventoryRepo.GetByID(ctx, itemID)
    if err != nil {
        if errors.Is(err, repository.ErrInventoryItemNotFound) {
            return nil, ErrItemNotFound
        }
        return nil, fmt.Errorf("get item: %w", err)
    }

    if item.UserID != userID {
        return nil, ErrNotItemOwner
    }

    if !item.IsEquipped {
        return nil, ErrAlreadyEquipped
    }

    if err := s.inventoryRepo.SetEquipped(ctx, itemID, false); err != nil {
        return nil, fmt.Errorf("unequip item: %w", err)
    }

    item.IsEquipped = false
    return item, nil
}