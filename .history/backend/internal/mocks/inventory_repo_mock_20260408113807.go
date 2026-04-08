package mocks

import (
    "context"
    "sync"
    "time"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/Shyramond/hakathon/backend/internal/models"
    "github.com/Shyramond/hakathon/backend/internal/repository"
)

type InventoryRepoMock struct {
    mu    sync.RWMutex
    items map[uuid.UUID]*models.InventoryItem
    userBenefits map[uuid.UUID]map[uuid.UUID]bool

    GetByUserIDCalls      int
    HasBenefitCalls       int
    AddCalls              int
    SetEquippedCalls      int
    UnequipAllByTypeCalls int
    GetByIDCalls          int
}

func NewInventoryRepoMock() *InventoryRepoMock {
    return &InventoryRepoMock{
        items:        make(map[uuid.UUID]*models.InventoryItem),
        userBenefits: make(map[uuid.UUID]map[uuid.UUID]bool),
    }
}

func (m *InventoryRepoMock) AddItem(item *models.InventoryItem) {
    m.mu.Lock()
    defer m.mu.Unlock()

    if item.ID == uuid.Nil {
        item.ID = uuid.New()
    }
    m.items[item.ID] = item

    if m.userBenefits[item.UserID] == nil {
        m.userBenefits[item.UserID] = make(map[uuid.UUID]bool)
    }
    m.userBenefits[item.UserID][item.BenefitID] = true
}

func (m *InventoryRepoMock) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.InventoryItem, error) {
    m.mu.Lock()
    m.GetByUserIDCalls++
    m.mu.Unlock()

    m.mu.RLock()
    defer m.mu.RUnlock()

    var result []models.InventoryItem
    for _, item := range m.items {
        if item.UserID == userID {
            result = append(result, *item)
        }
    }
    return result, nil
}

func (m *InventoryRepoMock) HasBenefit(ctx context.Context, userID, benefitID uuid.UUID) (bool, error) {
    m.mu.Lock()
    m.HasBenefitCalls++
    m.mu.Unlock()

    return m.hasBenefit(userID, benefitID), nil
}

func (m *InventoryRepoMock) HasBenefitTx(ctx context.Context, tx pgx.Tx, userID, benefitID uuid.UUID) (bool, error) {
    return m.HasBenefit(ctx, userID, benefitID)
}

func (m *InventoryRepoMock) hasBenefit(userID, benefitID uuid.UUID) bool {
    m.mu.RLock()
    defer m.mu.RUnlock()

    ub, ok := m.userBenefits[userID]
    if !ok {
        return false
    }
    return ub[benefitID]
}

func (m *InventoryRepoMock) Add(ctx context.Context, userID, benefitID uuid.UUID) (*models.InventoryItem, error) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.AddCalls++

    if m.userBenefits[userID] != nil && m.userBenefits[userID][benefitID] {
        return nil, repository.ErrAlreadyOwned
    }

    item := &models.InventoryItem{
        ID:          uuid.New(),
        UserID:      userID,
        BenefitID:   benefitID,
        PurchasedAt: time.Now(),
        IsEquipped:  false,
    }
    m.items[item.ID] = item

    if m.userBenefits[userID] == nil {
        m.userBenefits[userID] = make(map[uuid.UUID]bool)
    }
    m.userBenefits[userID][benefitID] = true

    return item, nil
}

func (m *InventoryRepoMock) AddTx(ctx context.Context, tx pgx.Tx, userID, benefitID uuid.UUID) (*models.InventoryItem, error) {
    return m.Add(ctx, userID, benefitID)
}

func (m *InventoryRepoMock) GetByID(ctx context.Context, id uuid.UUID) (*models.InventoryItem, error) {
    m.mu.Lock()
    m.GetByIDCalls++
    m.mu.Unlock()

    m.mu.RLock()
    defer m.mu.RUnlock()

    item, ok := m.items[id]
    if !ok {
        return nil, repository.ErrInventoryItemNotFound
    }
    copy := *item
    return &copy, nil
}

func (m *InventoryRepoMock) SetEquipped(ctx context.Context, id uuid.UUID, equipped bool) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.SetEquippedCalls++

    item, ok := m.items[id]
    if !ok {
        return repository.ErrInventoryItemNotFound
    }
    item.IsEquipped = equipped
    return nil
}

func (m *InventoryRepoMock) UnequipAllByType(ctx context.Context, userID uuid.UUID, benefitType models.BenefitType) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.UnequipAllByTypeCalls++

    for _, item := range m.items {
        if item.UserID == userID && item.BenefitType == benefitType {
            item.IsEquipped = false
        }
    }
    return nil
}