package router

import (
    "github.com/gin-gonic/gin"
    "github.com/Shyramond/hakathon/backend/internal/config"
    "github.com/Shyramond/hakathon/backend/internal/handler"
    "github.com/Shyramond/hakathon/backend/internal/middleware"
    "github.com/Shyramond/hakathon/backend/internal/repository"
    "github.com/Shyramond/hakathon/backend/internal/service"
    "github.com/Shyramond/hakathon/backend/pkg/xsolla"
)

type Dependencies struct {
    Config       *config.Config
    Repos        *repository.Repos
    Services     *service.Services
    XsollaClient *xsolla.Client
}

func Setup(deps Dependencies) *gin.Engine {
    gin.SetMode(deps.Config.Server.GinMode)

    r := gin.New()

    r.Use(gin.Recovery())
    r.Use(middleware.RequestLogging())
    r.Use(middleware.CORS())
    r.Use(middleware.SecurityHeaders())
    r.Use(middleware.SQLInjectionGuard())
    r.Use(middleware.MaxBodySize(1 << 20))
    r.Use(middleware.RateLimit(
        deps.Config.RateLimit.RPS,
        deps.Config.RateLimit.Burst,
    ))

    h := handler.NewHandlers(deps.Services, deps.Repos)

    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status":  "ok",
            "service": "hakathon-backend",
        })
    })

    v1 := r.Group("/api/v1")

    v1.POST("/auth/login", h.Auth.Login)

    v1.GET("/benefits", h.Benefit.ListBenefits)
    v1.GET("/benefits/:id", h.Benefit.GetBenefit)

    protected := v1.Group("")
    protected.Use(middleware.XsollaAuth(deps.XsollaClient))
    protected.Use(middleware.UserContext(deps.Repos.User))
    protected.Use(middleware.Idempotency(deps.Repos.Idempotency))
    {
        protected.GET("/profile", h.Profile.GetProfile)
        protected.GET("/transactions", h.Profile.GetTransactions)

        protected.POST("/shop/purchase/:id", h.Shop.Purchase)

        protected.GET("/inventory", h.Shop.GetInventory)
        protected.PATCH("/inventory/:id/equip", h.Shop.EquipItem)
        protected.PATCH("/inventory/:id/unequip", h.Shop.UnequipItem)
    }

    return r
}