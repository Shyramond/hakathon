package repository

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/Shyramond/hakathon/backend/internal/models"
)

var (
    ErrUserNotFound = errors.New("user not found")
)

type userRepoPG struct {
    pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
    return &userRepoPG{pool: pool}
}

const userSelectColumns = `
    id, xsolla_user_id, username, email, token_balance,
    last_login_at, created_at, updated_at
`

func scanUser(row pgx.Row) (*models.User, error) {
    var u models.User
    err := row.Scan(
        &u.ID, &u.XsollaUserID, &u.Username, &u.Email,
        &u.TokenBalance, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrUserNotFound
        }
        return nil, err
    }
    return &u, nil
}

func (r *userRepoPG) GetByXsollaID(ctx context.Context, xsollaUserID string) (*models.User, error) {
    query := fmt.Sprintf("SELECT %s FROM user_wallets WHERE xsolla_user_id = $1", userSelectColumns)
    u, err := scanUser(r.pool.QueryRow(ctx, query, xsollaUserID))
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            return nil, err
        }
        return nil, fmt.Errorf("get user by xsolla id: %w", err)
    }
    return u, nil
}

func (r *userRepoPG) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
    query := fmt.Sprintf("SELECT %s FROM user_wallets WHERE id = $1", userSelectColumns)
    u, err := scanUser(r.pool.QueryRow(ctx, query, id))
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            return nil, err
        }
        return nil, fmt.Errorf("get user by id: %w", err)
    }
    return u, nil
}

// GetByIDTx — версия GetByID внутри транзакции с FOR UPDATE (блокировка строки)
func (r *userRepoPG) GetByIDTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*models.User, error) {
    query := fmt.Sprintf("SELECT %s FROM user_wallets WHERE id = $1 FOR UPDATE", userSelectColumns)
    u, err := scanUser(tx.QueryRow(ctx, query, id))
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            return nil, err
        }
        return nil, fmt.Errorf("get user by id (tx): %w", err)
    }
    return u, nil
}

func (r *userRepoPG) Create(ctx context.Context, user *models.User) error {
    query := `
        INSERT INTO user_wallets (id, xsolla_user_id, username, email, token_balance, last_login_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING created_at, updated_at
    `
    if user.ID == uuid.Nil {
        user.ID = uuid.New()
    }
    err := r.pool.QueryRow(ctx, query,
        user.ID, user.XsollaUserID, user.Username, user.Email,
        user.TokenBalance, user.LastLoginAt,
    ).Scan(&user.CreatedAt, &user.UpdatedAt)
    if err != nil {
        return fmt.Errorf("create user: %w", err)
    }
    return nil
}

func (r *userRepoPG) UpdateBalance(ctx context.Context, id uuid.UUID, newBalance int) error {
    query := `UPDATE user_wallets SET token_balance = $2, updated_at = NOW() WHERE id = $1`
    tag, err := r.pool.Exec(ctx, query, id, newBalance)
    if err != nil {
        return fmt.Errorf("update balance: %w", err)
    }
    if tag.RowsAffected() == 0 {
        return ErrUserNotFound
    }
    return nil
}

func (r *userRepoPG) UpdateBalanceTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, newBalance int) error {
    query := `UPDATE user_wallets SET token_balance = $2, updated_at = NOW() WHERE id = $1`
    tag, err := tx.Exec(ctx, query, id, newBalance)
    if err != nil {
        return fmt.Errorf("update balance (tx): %w", err)
    }
    if tag.RowsAffected() == 0 {
        return ErrUserNotFound
    }
    return nil
}

func (r *userRepoPG) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
    query := `UPDATE user_wallets SET last_login_at = NOW(), updated_at = NOW() WHERE id = $1`
    tag, err := r.pool.Exec(ctx, query, id)
    if err != nil {
        return fmt.Errorf("update last login: %w", err)
    }
    if tag.RowsAffected() == 0 {
        return ErrUserNotFound
    }
    return nil
}

func (r *userRepoPG) UpdateLastLoginTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
    query := `UPDATE user_wallets SET last_login_at = NOW(), updated_at = NOW() WHERE id = $1`
    tag, err := tx.Exec(ctx, query, id)
    if err != nil {
        return fmt.Errorf("update last login (tx): %w", err)
    }
    if tag.RowsAffected() == 0 {
        return ErrUserNotFound
    }
    return nil
}

func (r *userRepoPG) GetOrCreate(ctx context.Context, xsollaUserID, username, email string) (*models.User, bool, error) {
    user, err := r.GetByXsollaID(ctx, xsollaUserID)
    if err == nil {
        return user, false, nil
    }
    if !errors.Is(err, ErrUserNotFound) {
        return nil, false, err
    }

    query := fmt.Sprintf(`
        INSERT INTO user_wallets (xsolla_user_id, username, email, token_balance, last_login_at)
        VALUES ($1, $2, $3, 0, $4)
        ON CONFLICT (xsolla_user_id) DO UPDATE
            SET username = EXCLUDED.username,
                email = EXCLUDED.email,
                updated_at = NOW()
        RETURNING %s
    `, userSelectColumns)

    now := time.Now()
    u, err := scanUser(r.pool.QueryRow(ctx, query, xsollaUserID, username, email, now))
    if err != nil {
        return nil, false, fmt.Errorf("get or create user: %w", err)
    }

    created := u.UpdatedAt.Sub(u.CreatedAt) < time.Second
    return u, created, nil
}