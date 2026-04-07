package routes

import (
	"github.com/gin-gonic/gin"

	"soulforge-backend/controllers"
	"soulforge-backend/middleware"
)

// SetupRoutes настраивает все маршруты
func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORSMiddleware())

	// Публичные маршруты
	public := router.Group("/api/v1")
	{
		// Аутентификация
		public.POST("/register", controllers.Register)
		public.POST("/login", controllers.Login)
	}

	// Защищенные маршруты
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		// Профиль пользователя
		protected.GET("/profile", controllers.GetProfile)

		// Инвентарь
		protected.GET("/inventory", controllers.GetInventory)
		protected.POST("/inventory/daily-reward", controllers.ClaimDailyReward)
		protected.GET("/inventory/daily-reward-config", controllers.GetDailyRewardConfig)
		protected.POST("/inventory/materials", controllers.AddMaterials)
		protected.GET("/inventory/history", controllers.GetCraftingHistory)
		protected.GET("/inventory/items", controllers.GetOwnedItems)

		// Крафтинг
		protected.POST("/crafting/upgrade", controllers.UpgradeMaterials)
		protected.GET("/crafting/recipes", controllers.GetCraftingRecipes)
		protected.POST("/crafting/item", controllers.CraftItem)

		// Премиум маршруты (требуют подписку)
		premium := protected.Group("/premium")
		premium.Use(middleware.RequireSubscription())
		{
			// Здесь будут премиум функции
			premium.GET("/exclusive-items", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Эксклюзивные предметы для подписчиков"})
			})
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "soulforge-backend"})
	})

	return router
}
