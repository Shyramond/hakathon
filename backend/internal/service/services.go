package service

import (
	"github.com/Shyramond/hakathon/backend/internal/config"
	"github.com/Shyramond/hakathon/backend/internal/repository"
	"github.com/Shyramond/hakathon/backend/pkg/xsolla"
)

type Services struct {
	Auth    *AuthService
	Token   *TokenService
	Benefit *BenefitService
	Shop    *ShopService
}

func NewServices(repos *repository.Repos, cfg *config.Config, xsollaClient *xsolla.Client) *Services {
	tokenService := NewTokenService(
		repos.User,
		repos.Transaction,
		repos.TxManager,
		cfg.Reward,
	)

	authService := NewAuthService(
		repos.User,
		tokenService,
		xsollaClient,
		cfg.Reward,
	)

	benefitService := NewBenefitService(repos.Benefit)

	shopService := NewShopService(
		repos.Benefit,
		repos.Inventory,
		tokenService,
		repos.TxManager,
	)

	return &Services{
		Auth:    authService,
		Token:   tokenService,
		Benefit: benefitService,
		Shop:    shopService,
	}
}
