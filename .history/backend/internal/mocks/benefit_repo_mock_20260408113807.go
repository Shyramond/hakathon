package mocks

import (
    "context"
    "sync"

    "github.com/google/uuid"
    "github.com/Shyramond/hakathon/backend/internal/models"
    "github.com/Shyramond/hakathon/backend/internal/repository"
)

type BenefitRepoMock struct {
    mu       sync.RWMutex
    benefits map[uuid.UUID]*models.Benefit

    GetAllCalls  int
    GetByIDCalls int
    CreateCalls  int
    UpdateCalls  int
    DeleteCalls  int
}

func NewBenefitRepoMock() *BenefitRepoMock {
    return &BenefitRepoMock{
        benefits: make(map[uuid.UUID]*models.Benefit),
    }
}

func (m *BenefitRepoMock) AddBenefit(b *models.Benefit) {
    m.mu.Lock()
    defer m.mu.Unlock()
    if b.ID == uuid.Nil {
        b.ID = uuid.New()
    }
    m.benefits[b.ID] = b
}

func (m *BenefitRepoMock) GetAll(ctx context.Context, filter models.BenefitFilter) ([]models.Benefit, error) {
    m.mu.Lock()
    m.GetAllCalls++
    m.mu.Unlock()

    m.mu.RLock()
    defer m.mu.RUnlock()

    var result []models.Benefit
    for _, b := range m.benefits {
        if filter.Type != nil && b.Type != *filter.Type {
            continue
        }
        if filter.IsActive != nil && b.IsActive != *filter.IsActive {
            continue
        }
        if filter.MinPrice != nil && b.PriceTokens < *filter.MinPrice {
            continue
        }
        if filter.MaxPrice != nil && b.PriceTokens > *filter.MaxPrice {
            continue
        }
        result = append(result, *b)
    }
    return result, nil
}

func (m *BenefitRepoMock) GetByID(ctx context.Context, id uuid.UUID) (*models.Benefit, error) {
    m.mu.Lock()
    m.GetByIDCalls++
    m.mu.Unlock()

    m.mu.RLock()
    defer m.mu.RUnlock()

    b, ok := m.benefits[id]
    if !ok {
        return nil, repository.ErrBenefitNotFound
    }
    copy := *b
    return &copy, nil
}

func (m *BenefitRepoMock) Create(ctx context.Context, benefit *models.Benefit) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.CreateCalls++

    if benefit.ID == uuid.Nil {
        benefit.ID = uuid.New()
    }
    m.benefits[benefit.ID] = benefit
    return nil
}

func (m *BenefitRepoMock) Update(ctx context.Context, benefit *models.Benefit) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.UpdateCalls++

    if _, ok := m.benefits[benefit.ID]; !ok {
        return repository.ErrBenefitNotFound
    }
    m.benefits[benefit.ID] = benefit
    return nil
}

func (m *BenefitRepoMock) Delete(ctx context.Context, id uuid.UUID) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.DeleteCalls++

    if _, ok := m.benefits[id]; !ok {
        return repository.ErrBenefitNotFound
    }
    delete(m.benefits, id)
    return nil
}