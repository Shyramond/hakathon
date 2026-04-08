package mocks

import (
    "context"
    "sync"

    "github.com/Shyramond/hakathon/backend/internal/repository"
)

type IdempotencyRepoMock struct {
    mu      sync.RWMutex
    records map[string]*repository.IdempotencyRecord

    GetCalls     int
    SetCalls     int
    CleanupCalls int
}

func NewIdempotencyRepoMock() *IdempotencyRepoMock {
    return &IdempotencyRepoMock{
        records: make(map[string]*repository.IdempotencyRecord),
    }
}

func (m *IdempotencyRepoMock) Get(ctx context.Context, key string) (*repository.IdempotencyRecord, error) {
    m.mu.Lock()
    m.GetCalls++
    m.mu.Unlock()

    m.mu.RLock()
    defer m.mu.RUnlock()

    rec, ok := m.records[key]
    if !ok {
        return nil, nil
    }
    return rec, nil
}

func (m *IdempotencyRepoMock) Set(ctx context.Context, record *repository.IdempotencyRecord) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.SetCalls++

    m.records[record.Key] = record
    return nil
}

func (m *IdempotencyRepoMock) Cleanup(ctx context.Context) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.CleanupCalls++
    return nil
}