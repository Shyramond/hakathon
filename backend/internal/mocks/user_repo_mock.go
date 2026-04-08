package mocks

import (
    "context"
    "sync"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/Shyramond/hakathon/backend/internal/models"
    "github.com/Shyramond/hakathon/backend/internal/repository"
)

type UserRepoMock struct {
    mu    sync.RWMutex
    users map[uuid.UUID]*models.User
    byXID map[string]*models.User

    GetByIDCalls          int
    GetByXsollaIDCalls    int
    CreateCalls           int
    UpdateBalanceCalls    int
    UpdateLastLoginCalls  int
    GetOrCreateCalls      int
    GetByIDTxCalls        int
    UpdateBalanceTxCalls  int
    UpdateLastLoginTxCalls int

    GetByIDErr       error
    UpdateBalanceErr error
    CreateErr        error
}

func NewUserRepoMock() *UserRepoMock {
    return &UserRepoMock{
        users: make(map[uuid.UUID]*models.User),
        byXID: make(map[string]*models.User),
    }
}

func (m *UserRepoMock) AddUser(u *models.User) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.users[u.ID] = u
    m.byXID[u.XsollaUserID] = u
}

func (m *UserRepoMock) GetUser(id uuid.UUID) *models.User {
    m.mu.RLock()
    defer m.mu.RUnlock()
    u := m.users[id]
    if u == nil {
        return nil
    }
    copy := *u
    return &copy
}

func (m *UserRepoMock) GetByXsollaID(ctx context.Context, xsollaUserID string) (*models.User, error) {
    m.mu.Lock()
    m.GetByXsollaIDCalls++
    m.mu.Unlock()

    m.mu.RLock()
    defer m.mu.RUnlock()

    u, ok := m.byXID[xsollaUserID]
    if !ok {
        return nil, repository.ErrUserNotFound
    }
    copy := *u
    return &copy, nil
}

func (m *UserRepoMock) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
    m.mu.Lock()
    m.GetByIDCalls++
    m.mu.Unlock()

    if m.GetByIDErr != nil {
        return nil, m.GetByIDErr
    }

    m.mu.RLock()
    defer m.mu.RUnlock()

    u, ok := m.users[id]
    if !ok {
        return nil, repository.ErrUserNotFound
    }
    copy := *u
    return &copy, nil
}

func (m *UserRepoMock) GetByIDTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*models.User, error) {
    m.mu.Lock()
    m.GetByIDTxCalls++
    m.mu.Unlock()

    return m.GetByID(ctx, id)
}

func (m *UserRepoMock) Create(ctx context.Context, user *models.User) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.CreateCalls++

    if m.CreateErr != nil {
        return m.CreateErr
    }

    if user.ID == uuid.Nil {
        user.ID = uuid.New()
    }
    m.users[user.ID] = user
    m.byXID[user.XsollaUserID] = user
    return nil
}

func (m *UserRepoMock) UpdateBalance(ctx context.Context, id uuid.UUID, newBalance int) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.UpdateBalanceCalls++

    if m.UpdateBalanceErr != nil {
        return m.UpdateBalanceErr
    }

    u, ok := m.users[id]
    if !ok {
        return repository.ErrUserNotFound
    }
    u.TokenBalance = newBalance
    return nil
}

func (m *UserRepoMock) UpdateBalanceTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, newBalance int) error {
    m.mu.Lock()
    m.UpdateBalanceTxCalls++
    m.mu.Unlock()

    return m.UpdateBalance(ctx, id, newBalance)
}

func (m *UserRepoMock) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.UpdateLastLoginCalls++

    _, ok := m.users[id]
    if !ok {
        return repository.ErrUserNotFound
    }
    return nil
}

func (m *UserRepoMock) UpdateLastLoginTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
    m.mu.Lock()
    m.UpdateLastLoginTxCalls++
    m.mu.Unlock()

    return m.UpdateLastLogin(ctx, id)
}

func (m *UserRepoMock) GetOrCreate(ctx context.Context, xsollaUserID, username, email string) (*models.User, bool, error) {
    m.mu.Lock()
    m.GetOrCreateCalls++
    m.mu.Unlock()

    existing, err := m.GetByXsollaID(ctx, xsollaUserID)
    if err == nil {
        return existing, false, nil
    }

    user := &models.User{
        ID:           uuid.New(),
        XsollaUserID: xsollaUserID,
        Username:     username,
        Email:        email,
        TokenBalance: 0,
    }
    if err := m.Create(ctx, user); err != nil {
        return nil, false, err
    }

    return user, true, nil
}