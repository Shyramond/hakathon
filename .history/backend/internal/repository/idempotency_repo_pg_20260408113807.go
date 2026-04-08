package repository

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type idempotencyRepoPG struct {
    pool *pgxpool.Pool
    ttl  time.Duration
}

func NewIdempotencyRepository(pool *pgxpool.Pool, ttl time.Duration) IdempotencyRepository {
    return &idempotencyRepoPG{pool: pool, ttl: ttl}
}

func (r *idempotencyRepoPG) Get(ctx context.Context, key string) (*IdempotencyRecord, error) {
    query := `
        SELECT key, user_id, response_code, response_body
        FROM idempotency_keys
        WHERE key = $1 AND expires_at > NOW()
    `

    var rec IdempotencyRecord
    err := r.pool.QueryRow(ctx, query, key).Scan(
        &rec.Key, &rec.UserID, &rec.ResponseCode, &rec.ResponseBody,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, fmt.Errorf("get idempotency key: %w", err)
    }

    return &rec, nil
}

func (r *idempotencyRepoPG) Set(ctx context.Context, record *IdempotencyRecord) error {
    query := `
        INSERT INTO idempotency_keys (key, user_id, response_code, response_body, expires_at)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (key) DO NOTHING
    `

    expiresAt := time.Now().Add(r.ttl)

    _, err := r.pool.Exec(ctx, query,
        record.Key, record.UserID, record.ResponseCode,
        record.ResponseBody, expiresAt,
    )
    if err != nil {
        return fmt.Errorf("set idempotency key: %w", err)
    }

    return nil
}

// Cleanup удаляет просроченные ключи (вызывать периодически)
func (r *idempotencyRepoPG) Cleanup(ctx context.Context) error {
    query := `DELETE FROM idempotency_keys WHERE expires_at < NOW()`

    tag, err := r.pool.Exec(ctx, query)
    if err != nil {
        return fmt.Errorf("cleanup idempotency keys: %w", err)
    }

    if tag.RowsAffected() > 0 {
        fmt.Printf("cleaned up %d expired idempotency keys\n", tag.RowsAffected())
    }

    return nil
}