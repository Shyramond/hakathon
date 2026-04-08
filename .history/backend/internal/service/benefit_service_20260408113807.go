package service

import (
    "context"
    "errors"
    "fmt"

    "github.com/google/uuid"
    "github.com/Shyramond/hakathon/backend/internal/models"
    "github.com/Shyramond/hakathon/backend/internal/repository"
)

type BenefitService struct {
    benefitRepo repository.BenefitRepository
}

func NewBenefitService(benefitRepo repository.BenefitRepository) *BenefitService {
    return &BenefitService{benefitRepo: benefitRepo}
}

func (s *BenefitService) GetAll(ctx context.Context, filter models.BenefitFilter) ([]models.Benefit, error) {
    if filter.IsActive == nil {
        active := true
        filter.IsActive = &active
    }

    benefits, err := s.benefitRepo.GetAll(ctx, filter)
    if err != nil {
        return nil, fmt.Errorf("get benefits: %w", err)
    }

    if benefits == nil {
        benefits = []models.Benefit{}
    }

    return benefits, nil
}

func (s *BenefitService) GetByID(ctx context.Context, id uuid.UUID) (*models.Benefit, error) {
    benefit, err := s.benefitRepo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, repository.ErrBenefitNotFound) {
            return nil, ErrBenefitNotFound
        }
        return nil, fmt.Errorf("get benefit: %w", err)
    }
    return benefit, nil
}

// Create создаёт новый бенефит (admin)
func (s *BenefitService) Create(ctx context.Context, benefit *models.Benefit) error {
    if err := s.validateBenefit(benefit); err != nil {
        return err
    }
    return s.benefitRepo.Create(ctx, benefit)
}

// Update обновляет бенефит (admin)
func (s *BenefitService) Update(ctx context.Context, benefit *models.Benefit) error {
    if err := s.validateBenefit(benefit); err != nil {
        return err
    }

    existing, err := s.benefitRepo.GetByID(ctx, benefit.ID)
    if err != nil {
        if errors.Is(err, repository.ErrBenefitNotFound) {
            return ErrBenefitNotFound
        }
        return err
    }
    _ = existing

    return s.benefitRepo.Update(ctx, benefit)
}

func (s *BenefitService) Delete(ctx context.Context, id uuid.UUID) error {
    benefit, err := s.benefitRepo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, repository.ErrBenefitNotFound) {
            return ErrBenefitNotFound
        }
        return err
    }

    benefit.IsActive = false
    return s.benefitRepo.Update(ctx, benefit)
}

func (s *BenefitService) validateBenefit(b *models.Benefit) error {
    if b.Name == "" {
        return fmt.Errorf("benefit name is required")
    }
    if b.PriceTokens < 0 {
        return fmt.Errorf("price cannot be negative")
    }

    validTypes := map[models.BenefitType]bool{
        models.BenefitTypePromotion:  true,
        models.BenefitTypeSkin:       true,
        models.BenefitTypeMascotSkin: true,
    }
    if !validTypes[b.Type] {
        return fmt.Errorf("invalid benefit type: %s", b.Type)
    }

    return nil
}