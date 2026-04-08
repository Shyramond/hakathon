package service

import (
    "context"
    "testing"

    "github.com/google/uuid"
    "github.com/Shyramond/hakathon/backend/internal/mocks"
    "github.com/Shyramond/hakathon/backend/internal/models"
)

func newTestBenefitService() (*BenefitService, *mocks.BenefitRepoMock) {
    repo := mocks.NewBenefitRepoMock()
    svc := NewBenefitService(repo)
    return svc, repo
}

func TestGetAll_DefaultActive(t *testing.T) {
    svc, repo := newTestBenefitService()

    repo.AddBenefit(&models.Benefit{Name: "Active Skin", IsActive: true, Type: models.BenefitTypeSkin, PriceTokens: 10})
    repo.AddBenefit(&models.Benefit{Name: "Inactive Skin", IsActive: false, Type: models.BenefitTypeSkin, PriceTokens: 20})

    benefits, err := svc.GetAll(context.Background(), models.BenefitFilter{})
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(benefits) != 1 {
        t.Errorf("expected 1 active benefit, got %d", len(benefits))
    }

    if len(benefits) > 0 && benefits[0].Name != "Active Skin" {
        t.Errorf("expected 'Active Skin', got '%s'", benefits[0].Name)
    }
}

func TestGetAll_FilterByType(t *testing.T) {
    svc, repo := newTestBenefitService()

    repo.AddBenefit(&models.Benefit{Name: "Skin 1", IsActive: true, Type: models.BenefitTypeSkin, PriceTokens: 10})
    repo.AddBenefit(&models.Benefit{Name: "Promo 1", IsActive: true, Type: models.BenefitTypePromotion, PriceTokens: 20})

    skinType := models.BenefitTypeSkin
    active := true
    benefits, err := svc.GetAll(context.Background(), models.BenefitFilter{
        Type:     &skinType,
        IsActive: &active,
    })
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(benefits) != 1 {
        t.Errorf("expected 1 skin benefit, got %d", len(benefits))
    }
}

func TestGetAll_FilterByPrice(t *testing.T) {
    svc, repo := newTestBenefitService()

    repo.AddBenefit(&models.Benefit{Name: "Cheap", IsActive: true, Type: models.BenefitTypeSkin, PriceTokens: 5})
    repo.AddBenefit(&models.Benefit{Name: "Mid", IsActive: true, Type: models.BenefitTypeSkin, PriceTokens: 50})
    repo.AddBenefit(&models.Benefit{Name: "Expensive", IsActive: true, Type: models.BenefitTypeSkin, PriceTokens: 200})

    minP := 10
    maxP := 100
    benefits, err := svc.GetAll(context.Background(), models.BenefitFilter{MinPrice: &minP, MaxPrice: &maxP})
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(benefits) != 1 {
        t.Errorf("expected 1 benefit in price range, got %d", len(benefits))
    }
}

func TestGetByID_NotFound(t *testing.T) {
    svc, _ := newTestBenefitService()

    _, err := svc.GetByID(context.Background(), uuid.New())
    if err == nil {
        t.Fatal("expected error for non-existent benefit")
    }
    if err != ErrBenefitNotFound {
        t.Errorf("expected ErrBenefitNotFound, got %v", err)
    }
}

func TestGetByID_Found(t *testing.T) {
    svc, repo := newTestBenefitService()

    id := uuid.New()
    repo.AddBenefit(&models.Benefit{
        ID:          id,
        Name:        "Cool Skin",
        IsActive:    true,
        Type:        models.BenefitTypeSkin,
        PriceTokens: 30,
    })

    b, err := svc.GetByID(context.Background(), id)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if b.Name != "Cool Skin" {
        t.Errorf("expected 'Cool Skin', got '%s'", b.Name)
    }
}

func TestValidateBenefit_TableDriven(t *testing.T) {
    svc, _ := newTestBenefitService()

    tests := []struct {
        name    string
        benefit *models.Benefit
        wantErr bool
    }{
        {
            name:    "valid skin",
            benefit: &models.Benefit{Name: "Skin", Type: models.BenefitTypeSkin, PriceTokens: 10},
            wantErr: false,
        },
        {
            name:    "valid promotion",
            benefit: &models.Benefit{Name: "Promo", Type: models.BenefitTypePromotion, PriceTokens: 0},
            wantErr: false,
        },
        {
            name:    "empty name",
            benefit: &models.Benefit{Name: "", Type: models.BenefitTypeSkin, PriceTokens: 10},
            wantErr: true,
        },
        {
            name:    "negative price",
            benefit: &models.Benefit{Name: "X", Type: models.BenefitTypeSkin, PriceTokens: -5},
            wantErr: true,
        },
        {
            name:    "invalid type",
            benefit: &models.Benefit{Name: "X", Type: "invalid_type", PriceTokens: 10},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := svc.validateBenefit(tt.benefit)
            if (err != nil) != tt.wantErr {
                t.Errorf("wantErr=%v, got err=%v", tt.wantErr, err)
            }
        })
    }
}

func TestDelete_SoftDelete(t *testing.T) {
    svc, repo := newTestBenefitService()

    id := uuid.New()
    repo.AddBenefit(&models.Benefit{
        ID:       id,
        Name:     "To Delete",
        IsActive: true,
        Type:     models.BenefitTypeSkin,
    })

    err := svc.Delete(context.Background(), id)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    b, _ := repo.GetByID(context.Background(), id)
    if b.IsActive {
        t.Error("expected benefit to be deactivated (soft delete)")
    }

    if repo.UpdateCalls != 1 {
        t.Errorf("expected 1 Update call, got %d", repo.UpdateCalls)
    }
}