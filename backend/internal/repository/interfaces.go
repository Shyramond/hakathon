package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/Shyramond/hakathon/backend/internal/models"
)

// TxManager управляет транзакциями БД
type TxManager interface {
    // RunInTx выполняет fn внутри транзакции. При ошибке — rollback, иначе — commit.
    RunInTx(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error
}

type UserRepository interface {
    GetByXsollaID(ctx context.Context, xsollaUserID string) (*models.User, error)
    GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
    Create(ctx context.Context, user *models.User) error
    UpdateBalance(ctx context.Context, id uuid.UUID, newBalance int) error
    UpdateLastLogin(ctx context.Context, id uuid.UUID) error
    GetOrCreate(ctx context.Context, xsollaUserID, username, email string) (*models.User, bool, error)

    // Tx-варианты для использования внутри транзакций
    GetByIDTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*models.User, error)
    UpdateBalanceTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, newBalance int) error
    UpdateLastLoginTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) error
}

type BenefitRepository interface {
    GetAll(ctx context.Context, filter models.BenefitFilter) ([]models.Benefit, error)
    GetByID(ctx context.Context, id uuid.UUID) (*models.Benefit, error)
    Create(ctx context.Context, benefit *models.Benefit) error
    Update(ctx context.Context, benefit *models.Benefit) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type InventoryRepository interface {
    GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.InventoryItem, error)
    HasBenefit(ctx context.Context, userID, benefitID uuid.UUID) (bool, error)
    Add(ctx context.Context, userID, benefitID uuid.UUID) (*models.InventoryItem, error)
    SetEquipped(ctx context.Context, id uuid.UUID, equipped bool) error
    UnequipAllByType(ctx context.Context, userID uuid.UUID, benefitType models.BenefitType) error
    GetByID(ctx context.Context, id uuid.UUID) (*models.InventoryItem, error)

    // Tx-варианты
    HasBenefitTx(ctx context.Context, tx pgx.Tx, userID, benefitID uuid.UUID) (bool, error)
    AddTx(ctx context.Context, tx pgx.Tx, userID, benefitID uuid.UUID) (*models.InventoryItem, error)
}

type TransactionRepository interface {
    Create(ctx context.Context, tx *models.TokenTransaction) error
    GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.TokenTransaction, int, error)

    // Tx-вариант
    CreateTx(ctx context.Context, pgxTx pgx.Tx, t *models.TokenTransaction) error
}

type IdempotencyRepository interface {
    Get(ctx context.Context, key string) (*IdempotencyRecord, error)
    Set(ctx context.Context, record *IdempotencyRecord) error
    Cleanup(ctx context.Context) error
}

type IdempotencyRecord struct {
    Key          string    `json:"key"`
    UserID       uuid.UUID `json:"user_id"`
    ResponseCode int       `json:"response_code"`
    ResponseBody []byte    `json:"response_body"`
    ExpiresAt    time.Time `json:"expires_at"`
}