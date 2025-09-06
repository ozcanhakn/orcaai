package main

import (
	"log"
	"net/http"
	"os"

	"orcaai/config"
	"orcaai/database"
	"orcaai/handlers"
	"orcaai/middleware"
	"orcaai/models"
	"orcaai/orchestrator"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize configuration
	cfg := config.Load()

	// Initialize database connection
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Redis connection
	redis := database.ConnectRedis(cfg.RedisURL)
	defer redis.Close()

	// Initialize cache system
	cacheType := os.Getenv("CACHE_TYPE")
	if cacheType == "" {
		cacheType = "memory"
	}

	var cacheAddr, cachePassword string
	var cacheDB int

	if cacheType == "redis" {
		cacheAddr = os.Getenv("CACHE_REDIS_ADDR")
		if cacheAddr == "" {
			cacheAddr = "localhost:6379"
		}
		cachePassword = os.Getenv("CACHE_REDIS_PASSWORD")
		// cacheDB is 0 by default
	}

	if err := orchestrator.InitializeCache(cacheType, cacheAddr, cachePassword, cacheDB); err != nil {
		log.Printf("Warning: Failed to initialize cache: %v", err)
	} else {
		log.Printf("Cache initialized successfully: %s", cacheType)
	}
	defer orchestrator.CloseCache()

	// Initialize metrics
	orchestrator.InitializeMetrics()

	// Initialize Gin router
	r := gin.Default()

	// Apply middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// CORS middleware - using Gin's built-in CORS or a custom implementation
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "orcaai-backend",
			"version": "1.0.0",
		})
	})

	// Metrics endpoint for Prometheus
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	api := r.Group("/api/v1")
	{
		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/register", handlers.Register)
			auth.POST("/refresh", handlers.RefreshToken)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// AI query endpoint
			protected.POST("/ai/query", handlers.AIQuery)
			protected.GET("/ai/providers", handlers.GetProviders)

			// User management (user can view their own profile)
			protected.GET("/user/profile", handlers.GetProfile)
			protected.PUT("/user/profile", handlers.UpdateProfile)

			// API Key management (user can manage their own API keys)
			protected.GET("/keys", handlers.GetAPIKeys)
			protected.POST("/keys", handlers.CreateAPIKey)
			protected.DELETE("/keys/:id", handlers.DeleteAPIKey)

			// Metrics and analytics (all users can view basic metrics)
			protected.GET("/metrics", handlers.GetMetrics)
			protected.GET("/metrics/usage", handlers.GetUsageMetrics)
			protected.GET("/metrics/cost", handlers.GetCostMetrics)

			// Admin routes (only admin users can access)
			admin := protected.Group("/admin")
			admin.Use(middleware.RoleMiddleware(models.RoleAdmin))
			{
				admin.GET("/users", handlers.GetAllUsers)
				admin.PUT("/users/:id", handlers.UpdateUser)
				admin.DELETE("/users/:id", handlers.DeleteUser)
				admin.GET("/metrics/detailed", handlers.GetDetailedMetrics)
			}

			// Enterprise routes (enterprise users and admins can access)
			enterprise := protected.Group("/enterprise")
			enterprise.Use(middleware.RoleMiddleware(models.RoleEnterprise, models.RoleAdmin))
			{
				enterprise.GET("/features", handlers.GetEnterpriseFeatures)
				enterprise.POST("/features", handlers.EnableEnterpriseFeature)
			}
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
