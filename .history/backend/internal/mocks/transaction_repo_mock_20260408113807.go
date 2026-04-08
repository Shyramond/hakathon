package mocks

import (
    "context"
    "sync"
    "time"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/Shyramond/hakathon/backend/internal/models"
)

type TransactionRepoMock struct {
    mu           sync.RWMutex
    transactions []models.TokenTransaction

    CreateCalls    int
    GetByUserCalls int
}

func NewTransactionRepoMock() *TransactionRepoMock {
    return &TransactionRepoMock{}
}

func (m *TransactionRepoMock) Create(ctx context.Context, t *models.TokenTransaction) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.CreateCalls++

    t.ID = uuid.New()
    t.CreatedAt = time.Now()
    m.transactions = append(m.transactions, *t)
    return nil
}

func (m *TransactionRepoMock) CreateTx(ctx context.Context, pgxTx pgx.Tx, t *models.TokenTransaction) error {
    return m.Create(ctx, t)
}

func (m *TransactionRepoMock) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.TokenTransaction, int, error) {
    m.mu.Lock()
    m.GetByUserCalls++
    m.mu.Unlock()

    m.mu.RLock()
    defer m.mu.RUnlock()

    var filtered []models.TokenTransaction
    for _, t := range m.transactions {
        if t.UserID == userID {
            filtered = append(filtered, t)
        }
    }

    total := len(filtered)

    if offset >= len(filtered) {
        return []models.TokenTransaction{}, total, nil
    }
    end := offset + limit
    if end > len(filtered) {
        end = len(filtered)
    }

    return filtered[offset:end], total, nil
}

func (m *TransactionRepoMock) GetAll() []models.TokenTransaction {
    m.mu.RLock()
    defer m.mu.RUnlock()
    result := make([]models.TokenTransaction, len(m.transactions))
    copy(result, m.transactions)
    return result
}