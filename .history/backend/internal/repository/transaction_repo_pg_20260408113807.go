package repository

import (
    "context"
    "fmt"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/Shyramond/hakathon/backend/internal/models"
)

type transactionRepoPG struct {
    pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) TransactionRepository {
    return &transactionRepoPG{pool: pool}
}

func (r *transactionRepoPG) Create(ctx context.Context, t *models.TokenTransaction) error {
    return createTransaction(ctx, r.pool, t)
}

func (r *transactionRepoPG) CreateTx(ctx context.Context, pgxTx pgx.Tx, t *models.TokenTransaction) error {
    return createTransaction(ctx, pgxTx, t)
}

func (r *transactionRepoPG) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.TokenTransaction, int, error) {
    countQuery := `SELECT COUNT(*) FROM token_transactions WHERE user_id = $1`
    var total int
    err := r.pool.QueryRow(ctx, countQuery, userID).Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("count transactions: %w", err)
    }

    if total == 0 {
        return []models.TokenTransaction{}, 0, nil
    }

    query := `
        SELECT id, user_id, amount, type, reference_id, description, created_at
        FROM token_transactions
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

    rows, err := r.pool.Query(ctx, query, userID, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("query transactions: %w", err)
    }
    defer rows.Close()

    var transactions []models.TokenTransaction
    for rows.Next() {
        var t models.TokenTransaction
        err := rows.Scan(
            &t.ID, &t.UserID, &t.Amount, &t.Type,
            &t.ReferenceID, &t.Description, &t.CreatedAt,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("scan transaction: %w", err)
        }
        transactions = append(transactions, t)
    }

    if err := rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("iterate transactions: %w", err)
    }

    return transactions, total, nil
}

func createTransaction(ctx context.Context, q interface {
    QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}, t *models.TokenTransaction) error {
    query := `
        INSERT INTO token_transactions (user_id, amount, type, reference_id, description)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at
    `
    err := q.QueryRow(ctx, query,
        t.UserID, t.Amount, t.Type,
        t.ReferenceID, t.Description,
    ).Scan(&t.ID, &t.CreatedAt)
    if err != nil {
        return fmt.Errorf("create transaction: %w", err)
    }
    return nil
}