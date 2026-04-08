package repository

import (
    "context"
    "errors"
    "fmt"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/Shyramond/hakathon/backend/internal/models"
)

var (
    ErrInventoryItemNotFound = errors.New("inventory item not found")
    ErrAlreadyOwned          = errors.New("benefit already owned")
)

type inventoryRepoPG struct {
    pool *pgxpool.Pool
}

func NewInventoryRepository(pool *pgxpool.Pool) InventoryRepository {
    return &inventoryRepoPG{pool: pool}
}

func (r *inventoryRepoPG) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.InventoryItem, error) {
    query := `
        SELECT
            ui.id, ui.user_id, ui.benefit_id, ui.purchased_at, ui.is_equipped,
            b.name, b.type, b.image_url
        FROM user_inventory ui
        JOIN benefits b ON b.id = ui.benefit_id
        WHERE ui.user_id = $1
        ORDER BY ui.purchased_at DESC
    `

    rows, err := r.pool.Query(ctx, query, userID)
    if err != nil {
        return nil, fmt.Errorf("query inventory: %w", err)
    }
    defer rows.Close()

    return scanInventoryRows(rows)
}

func (r *inventoryRepoPG) HasBenefit(ctx context.Context, userID, benefitID uuid.UUID) (bool, error) {
    return hasBenefitQuery(ctx, r.pool, userID, benefitID)
}

func (r *inventoryRepoPG) HasBenefitTx(ctx context.Context, tx pgx.Tx, userID, benefitID uuid.UUID) (bool, error) {
    return hasBenefitQuery(ctx, tx, userID, benefitID)
}

func (r *inventoryRepoPG) Add(ctx context.Context, userID, benefitID uuid.UUID) (*models.InventoryItem, error) {
    return addInventoryItem(ctx, r.pool, userID, benefitID)
}

func (r *inventoryRepoPG) AddTx(ctx context.Context, tx pgx.Tx, userID, benefitID uuid.UUID) (*models.InventoryItem, error) {
    return addInventoryItem(ctx, tx, userID, benefitID)
}

func (r *inventoryRepoPG) GetByID(ctx context.Context, id uuid.UUID) (*models.InventoryItem, error) {
    query := `
        SELECT
            ui.id, ui.user_id, ui.benefit_id, ui.purchased_at, ui.is_equipped,
            b.name, b.type, b.image_url
        FROM user_inventory ui
        JOIN benefits b ON b.id = ui.benefit_id
        WHERE ui.id = $1
    `

    var item models.InventoryItem
    err := r.pool.QueryRow(ctx, query, id).Scan(
        &item.ID, &item.UserID, &item.BenefitID,
        &item.PurchasedAt, &item.IsEquipped,
        &item.BenefitName, &item.BenefitType, &item.ImageURL,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrInventoryItemNotFound
        }
        return nil, fmt.Errorf("get inventory item: %w", err)
    }

    return &item, nil
}

func (r *inventoryRepoPG) SetEquipped(ctx context.Context, id uuid.UUID, equipped bool) error {
    query := `UPDATE user_inventory SET is_equipped = $2 WHERE id = $1`
    tag, err := r.pool.Exec(ctx, query, id, equipped)
    if err != nil {
        return fmt.Errorf("set equipped: %w", err)
    }
    if tag.RowsAffected() == 0 {
        return ErrInventoryItemNotFound
    }
    return nil
}

func (r *inventoryRepoPG) UnequipAllByType(ctx context.Context, userID uuid.UUID, benefitType models.BenefitType) error {
    query := `
        UPDATE user_inventory ui
        SET is_equipped = false
        FROM benefits b
        WHERE ui.benefit_id = b.id
          AND ui.user_id = $1
          AND b.type = $2
          AND ui.is_equipped = true
    `
    _, err := r.pool.Exec(ctx, query, userID, string(benefitType))
    if err != nil {
        return fmt.Errorf("unequip all by type: %w", err)
    }
    return nil
}

// === Общие хелперы для pool и tx ===

// querier — общий интерфейс для pgxpool.Pool и pgx.Tx
type querier interface {
    QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
    Exec(ctx context.Context, sql string, args ...any) (pgx.Row, error)
}

// hasBenefitQuery работает и с pool, и с tx
func hasBenefitQuery(ctx context.Context, q interface{ QueryRow(ctx context.Context, sql string, args ...any) pgx.Row }, userID, benefitID uuid.UUID) (bool, error) {
    query := `
        SELECT EXISTS(
            SELECT 1 FROM user_inventory
            WHERE user_id = $1 AND benefit_id = $2
        )
    `
    var exists bool
    err := q.QueryRow(ctx, query, userID, benefitID).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("check benefit ownership: %w", err)
    }
    return exists, nil
}

// addInventoryItem работает и с pool, и с tx
func addInventoryItem(ctx context.Context, q interface{ QueryRow(ctx context.Context, sql string, args ...any) pgx.Row }, userID, benefitID uuid.UUID) (*models.InventoryItem, error) {
    query := `
        INSERT INTO user_inventory (user_id, benefit_id)
        VALUES ($1, $2)
        RETURNING id, user_id, benefit_id, purchased_at, is_equipped
    `
    var item models.InventoryItem
    err := q.QueryRow(ctx, query, userID, benefitID).Scan(
        &item.ID, &item.UserID, &item.BenefitID,
        &item.PurchasedAt, &item.IsEquipped,
    )
    if err != nil {
        if isDuplicateKeyError(err) {
            return nil, ErrAlreadyOwned
        }
        return nil, fmt.Errorf("add to inventory: %w", err)
    }
    return &item, nil
}

func scanInventoryRows(rows pgx.Rows) ([]models.InventoryItem, error) {
    var items []models.InventoryItem
    for rows.Next() {
        var item models.InventoryItem
        err := rows.Scan(
            &item.ID, &item.UserID, &item.BenefitID,
            &item.PurchasedAt, &item.IsEquipped,
            &item.BenefitName, &item.BenefitType, &item.ImageURL,
        )
        if err != nil {
            return nil, fmt.Errorf("scan inventory item: %w", err)
        }
        items = append(items, item)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("iterate inventory: %w", err)
    }
    return items, nil
}

func isDuplicateKeyError(err error) bool {
    if err == nil {
        return false
    }
    var pgErr interface{ SQLState() string }
    if errors.As(err, &pgErr) {
        return pgErr.SQLState() == "23505"
    }
    return false
}