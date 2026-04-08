package repository

import (
    "context"
    "errors"
    "fmt"
    "strings"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/Shyramond/hakathon/backend/internal/models"
)

var (
    ErrBenefitNotFound = errors.New("benefit not found")
)

type benefitRepoPG struct {
    pool *pgxpool.Pool
}

func NewBenefitRepository(pool *pgxpool.Pool) BenefitRepository {
    return &benefitRepoPG{pool: pool}
}

func (r *benefitRepoPG) GetAll(ctx context.Context, filter models.BenefitFilter) ([]models.Benefit, error) {
    var (
        conditions []string
        args       []interface{}
        argIdx     = 1
    )

    if filter.Type != nil {
        conditions = append(conditions, fmt.Sprintf("type = $%d", argIdx))
        args = append(args, string(*filter.Type))
        argIdx++
    }

    if filter.IsActive != nil {
        conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIdx))
        args = append(args, *filter.IsActive)
        argIdx++
    }

    if filter.MinPrice != nil {
        conditions = append(conditions, fmt.Sprintf("price_tokens >= $%d", argIdx))
        args = append(args, *filter.MinPrice)
        argIdx++
    }

    if filter.MaxPrice != nil {
        conditions = append(conditions, fmt.Sprintf("price_tokens <= $%d", argIdx))
        args = append(args, *filter.MaxPrice)
        argIdx++
    }

    query := `
        SELECT id, name, description, type, price_tokens, image_url,
               is_active, metadata, created_at, updated_at
        FROM benefits
    `

    if len(conditions) > 0 {
        query += " WHERE " + strings.Join(conditions, " AND ")
    }

    query += " ORDER BY created_at DESC"

    rows, err := r.pool.Query(ctx, query, args...)
    if err != nil {
        return nil, fmt.Errorf("query benefits: %w", err)
    }
    defer rows.Close()

    var benefits []models.Benefit
    for rows.Next() {
        var b models.Benefit
        err := rows.Scan(
            &b.ID, &b.Name, &b.Description, &b.Type,
            &b.PriceTokens, &b.ImageURL, &b.IsActive,
            &b.Metadata, &b.CreatedAt, &b.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("scan benefit: %w", err)
        }
        benefits = append(benefits, b)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("iterate benefits: %w", err)
    }

    return benefits, nil
}

func (r *benefitRepoPG) GetByID(ctx context.Context, id uuid.UUID) (*models.Benefit, error) {
    query := `
        SELECT id, name, description, type, price_tokens, image_url,
               is_active, metadata, created_at, updated_at
        FROM benefits
        WHERE id = $1
    `

    var b models.Benefit
    err := r.pool.QueryRow(ctx, query, id).Scan(
        &b.ID, &b.Name, &b.Description, &b.Type,
        &b.PriceTokens, &b.ImageURL, &b.IsActive,
        &b.Metadata, &b.CreatedAt, &b.UpdatedAt,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrBenefitNotFound
        }
        return nil, fmt.Errorf("get benefit by id: %w", err)
    }

    return &b, nil
}

func (r *benefitRepoPG) Create(ctx context.Context, benefit *models.Benefit) error {
    query := `
        INSERT INTO benefits (name, description, type, price_tokens, image_url, is_active, metadata)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at
    `

    if benefit.Metadata == nil {
        benefit.Metadata = []byte("{}")
    }

    err := r.pool.QueryRow(ctx, query,
        benefit.Name, benefit.Description, benefit.Type,
        benefit.PriceTokens, benefit.ImageURL, benefit.IsActive,
        benefit.Metadata,
    ).Scan(&benefit.ID, &benefit.CreatedAt, &benefit.UpdatedAt)

    if err != nil {
        return fmt.Errorf("create benefit: %w", err)
    }

    return nil
}

func (r *benefitRepoPG) Update(ctx context.Context, benefit *models.Benefit) error {
    query := `
        UPDATE benefits
        SET name = $2, description = $3, type = $4, price_tokens = $5,
            image_url = $6, is_active = $7, metadata = $8, updated_at = NOW()
        WHERE id = $1
        RETURNING updated_at
    `

    err := r.pool.QueryRow(ctx, query,
        benefit.ID, benefit.Name, benefit.Description, benefit.Type,
        benefit.PriceTokens, benefit.ImageURL, benefit.IsActive,
        benefit.Metadata,
    ).Scan(&benefit.UpdatedAt)

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return ErrBenefitNotFound
        }
        return fmt.Errorf("update benefit: %w", err)
    }

    return nil
}

func (r *benefitRepoPG) Delete(ctx context.Context, id uuid.UUID) error {
    query := `DELETE FROM benefits WHERE id = $1`

    tag, err := r.pool.Exec(ctx, query, id)
    if err != nil {
        return fmt.Errorf("delete benefit: %w", err)
    }
    if tag.RowsAffected() == 0 {
        return ErrBenefitNotFound
    }

    return nil
}