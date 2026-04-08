package repository

import (
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

type Repos struct {
    User        UserRepository
    Benefit     BenefitRepository
    Inventory   InventoryRepository
    Transaction TransactionRepository
    Idempotency IdempotencyRepository
    TxManager   TxManager
}

func NewRepos(pool *pgxpool.Pool, idempotencyTTL time.Duration) *Repos {
    return &Repos{
        User:        NewUserRepository(pool),
        Benefit:     NewBenefitRepository(pool),
        Inventory:   NewInventoryRepository(pool),
        Transaction: NewTransactionRepository(pool),
        Idempotency: NewIdempotencyRepository(pool, idempotencyTTL),
        TxManager:   NewTxManager(pool),
    }
}