package routes

import (
	"net/http"
	"taskmango/authsvc/internal/config"
	"taskmango/authsvc/internal/controllers"
	"taskmango/authsvc/internal/repositories"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	userRepo := repositories.NewUserRepository(db)
	authController := controllers.NewAuthController(userRepo, cfg)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", authController.Login)
	}

	return router
}

func SetupProbeRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	router.GET("/readyz", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "not ready",
				"error":  err.Error(),
			})
			return
		}

		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "not ready",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	return router
}