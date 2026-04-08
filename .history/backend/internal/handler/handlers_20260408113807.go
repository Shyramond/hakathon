package handler

import (
    "github.com/Shyramond/hakathon/backend/internal/repository"
    "github.com/Shyramond/hakathon/backend/internal/service"
)

type Handlers struct {
    Auth    *AuthHandler
    Benefit *BenefitHandler
    Shop    *ShopHandler
    Profile *ProfileHandler
}

func NewHandlers(services *service.Services, repos *repository.Repos) *Handlers {
    return &Handlers{
        Auth:    NewAuthHandler(services.Auth),
        Benefit: NewBenefitHandler(services.Benefit),
        Shop:    NewShopHandler(services.Shop),
        Profile: NewProfileHandler(services.Shop, services.Token, repos.Transaction),
    }
}