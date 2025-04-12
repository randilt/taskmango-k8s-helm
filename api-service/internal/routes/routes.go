package routes

import (
	"net/http"

	"taskmango/apisvc/internal/config"
	"taskmango/apisvc/internal/controllers"
	"taskmango/apisvc/internal/middlewares"
	"taskmango/apisvc/internal/repositories"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Initialize repositories
	taskRepo := repositories.NewTaskRepository(db)
	tagRepo := repositories.NewTagRepository(db)

	// Initialize middleware
	authMiddleware := middlewares.AuthMiddleware(cfg)

	// Initialize controllers
	taskController := controllers.NewTaskController(taskRepo, tagRepo)

	// API routes
	apiGroup := router.Group("/api")
	apiGroup.Use(authMiddleware)
	{
		// Tasks endpoints
		tasksGroup := apiGroup.Group("/tasks")
		{
			tasksGroup.GET("", taskController.GetTasks)
			tasksGroup.GET("/:id", taskController.GetTaskByID)
			tasksGroup.POST("", taskController.CreateTask)
			tasksGroup.PUT("/:id", taskController.UpdateTask)
			tasksGroup.DELETE("/:id", taskController.DeleteTask)
		}

		// Tags endpoint
		apiGroup.GET("/tags", taskController.GetTags)
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