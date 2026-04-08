package service

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/Shyramond/hakathon/backend/internal/config"
    "github.com/Shyramond/hakathon/backend/internal/mocks"
    "github.com/Shyramond/hakathon/backend/internal/models"
)

type shopTestEnv struct {
    shopService   *ShopService
    userRepo      *mocks.UserRepoMock
    benefitRepo   *mocks.BenefitRepoMock
    inventoryRepo *mocks.InventoryRepoMock
    txRepo        *mocks.TransactionRepoMock
    txMgr         *mocks.TxManagerMock
}

func newShopTestEnv() *shopTestEnv {
    userRepo := mocks.NewUserRepoMock()
    benefitRepo := mocks.NewBenefitRepoMock()
    inventoryRepo := mocks.NewInventoryRepoMock()
    txRepo := mocks.NewTransactionRepoMock()
    txMgr := mocks.NewTxManagerMock()

    tokenService := NewTokenService(userRepo, txRepo, txMgr, config.RewardConfig{
        Min: 10, Max: 50, CooldownDuration: 24 * time.Hour,
    })

    shopService := NewShopService(benefitRepo, inventoryRepo, tokenService, txMgr)

    return &shopTestEnv{
        shopService:   shopService,
        userRepo:      userRepo,
        benefitRepo:   benefitRepo,
        inventoryRepo: inventoryRepo,
        txRepo:        txRepo,
        txMgr:         txMgr,
    }
}

func TestPurchase_Success(t *testing.T) {
    env := newShopTestEnv()

    userID := uuid.New()
    benefitID := uuid.New()

    env.userRepo.AddUser(&models.User{
        ID:           userID,
        TokenBalance: 100,
    })

    env.benefitRepo.AddBenefit(&models.Benefit{
        ID:          benefitID,
        Name:        "Cool Skin",
        PriceTokens: 30,
        IsActive:    true,
        Type:        models.BenefitTypeSkin,
    })

    result, err := env.shopService.Purchase(context.Background(), userID, benefitID)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if result.TokensSpent != 30 {
        t.Errorf("expected 30 tokens spent, got %d", result.TokensSpent)
    }

    if result.NewBalance != 70 {
        t.Errorf("expected new balance 70, got %d", result.NewBalance)
    }

    if result.InventoryItem == nil {
        t.Fatal("expected inventory item")
    }

    // Проверяем баланс в моке
    user := env.userRepo.GetUser(userID)
    if user.TokenBalance != 70 {
        t.Errorf("expected stored balance 70, got %d", user.TokenBalance)
    }

    // Проверяем что транзакция записана
    if env.txRepo.CreateCalls != 1 {
        t.Errorf("expected 1 transaction, got %d", env.txRepo.CreateCalls)
    }

    // Транзакция должна быть с RunInTx
    if env.txMgr.RunInTxCalls != 1 {
        t.Errorf("expected 1 RunInTx call, got %d", env.txMgr.RunInTxCalls)
    }
}

func TestPurchase_InsufficientBalance(t *testing.T) {
    env := newShopTestEnv()

    userID := uuid.New()
    benefitID := uuid.New()

    env.userRepo.AddUser(&models.User{
        ID:           userID,
        TokenBalance: 10, // Мало токенов
    })

    env.benefitRepo.AddBenefit(&models.Benefit{
        ID:          benefitID,
        Name:        "Expensive Skin",
        PriceTokens: 100,
        IsActive:    true,
        Type:        models.BenefitTypeSkin,
    })

    _, err := env.shopService.Purchase(context.Background(), userID, benefitID)
    if !errors.Is(err, ErrInsufficientBalance) {
        t.Errorf("expected ErrInsufficientBalance, got %v", err)
    }
}

func TestPurchase_AlreadyOwned(t *testing.T) {
    env := newShopTestEnv()

    userID := uuid.New()
    benefitID := uuid.New()

    env.userRepo.AddUser(&models.User{
        ID:           userID,
        TokenBalance: 100,
    })

    env.benefitRepo.AddBenefit(&models.Benefit{
        ID:          benefitID,
        Name:        "Owned Skin",
        PriceTokens: 10,
        IsActive:    true,
        Type:        models.BenefitTypeSkin,
    })

    // Уже владеет
    env.inventoryRepo.AddItem(&models.InventoryItem{
        UserID:    userID,
        BenefitID: benefitID,
    })

    _, err := env.shopService.Purchase(context.Background(), userID, benefitID)
    if !errors.Is(err, ErrAlreadyOwned) {
        t.Errorf("expected ErrAlreadyOwned, got %v", err)
    }
}

func TestPurchase_BenefitNotActive(t *testing.T) {
    env := newShopTestEnv()

    userID := uuid.New()
    benefitID := uuid.New()

    env.userRepo.AddUser(&models.User{ID: userID, TokenBalance: 100})
    env.benefitRepo.AddBenefit(&models.Benefit{
        ID: benefitID, Name: "Disabled", PriceTokens: 10, IsActive: false,
    })

    _, err := env.shopService.Purchase(context.Background(), userID, benefitID)
    if !errors.Is(err, ErrBenefitNotActive) {
        t.Errorf("expected ErrBenefitNotActive, got %v", err)
    }
}

func TestPurchase_BenefitNotFound(t *testing.T) {
    env := newShopTestEnv()

    userID := uuid.New()
    env.userRepo.AddUser(&models.User{ID: userID, TokenBalance: 100})

    _, err := env.shopService.Purchase(context.Background(), userID, uuid.New())
    if !errors.Is(err, ErrBenefitNotFound) {
        t.Errorf("expected ErrBenefitNotFound, got %v", err)
    }
}

func TestPurchase_FreeBenefit(t *testing.T) {
    env := newShopTestEnv()

    userID := uuid.New()
    benefitID := uuid.New()

    env.userRepo.AddUser(&models.User{ID: userID, TokenBalance: 0})
    env.benefitRepo.AddBenefit(&models.Benefit{
        ID: benefitID, Name: "Free Promo", PriceTokens: 0, IsActive: true, Type: models.BenefitTypePromotion,
    })

    result, err := env.shopService.Purchase(context.Background(), userID, benefitID)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if result.TokensSpent != 0 {
        t.Errorf("expected 0 tokens spent, got %d", result.TokensSpent)
    }

    if result.NewBalance != 0 {
        t.Errorf("expected balance 0, got %d", result.NewBalance)
    }
}

func TestEquipItem_Success(t *testing.T) {
    env := newShopTestEnv()

    userID := uuid.New()
    itemID := uuid.New()

    env.inventoryRepo.AddItem(&models.InventoryItem{
        ID:          itemID,
        UserID:      userID,
        BenefitID:   uuid.New(),
        BenefitType: models.BenefitTypeSkin,
        BenefitName: "Skin",
        IsEquipped:  false,
    })

    item, err := env.shopService.EquipItem(context.Background(), userID, itemID)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if !item.IsEquipped {
        t.Error("expected item to be equipped")
    }
}

func TestEquipItem_NotOwner(t *testing.T) {
    env := newShopTestEnv()

    ownerID := uuid.New()
    attackerID := uuid.New()
    itemID := uuid.New()

    env.inventoryRepo.AddItem(&models.InventoryItem{
        ID:          itemID,
        UserID:      ownerID,
        BenefitType: models.BenefitTypeSkin,
        IsEquipped:  false,
    })

    _, err := env.shopService.EquipItem(context.Background(), attackerID, itemID)
    if !errors.Is(err, ErrNotItemOwner) {
        t.Errorf("expected ErrNotItemOwner, got %v", err)
    }
}

func TestEquipItem_AlreadyEquipped(t *testing.T) {
    env := newShopTestEnv()

    userID := uuid.New()
    itemID := uuid.New()

    env.inventoryRepo.AddItem(&models.InventoryItem{
        ID:          itemID,
        UserID:      userID,
        BenefitType: models.BenefitTypeSkin,
        IsEquipped:  true,
    })

    _, err := env.shopService.EquipItem(context.Background(), userID, itemID)
    if !errors.Is(err, ErrAlreadyEquipped) {
        t.Errorf("expected ErrAlreadyEquipped, got %v", err)
    }
}

func TestUnequipItem_Success(t *testing.T) {
    env := newShopTestEnv()

    userID := uuid.New()
    itemID := uuid.New()

    env.inventoryRepo.AddItem(&models.InventoryItem{
        ID:          itemID,
        UserID:      userID,
        BenefitType: models.BenefitTypeSkin,
        IsEquipped:  true,
    })

    item, err := env.shopService.UnequipItem(context.Background(), userID, itemID)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if item.IsEquipped {
        t.Error("expected item to be unequipped")
    }
}

func TestGetInventory_Empty(t *testing.T) {
    env := newShopTestEnv()

    items, err := env.shopService.GetInventory(context.Background(), uuid.New())
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(items) != 0 {
        t.Errorf("expected empty inventory, got %d items", len(items))
    }
}