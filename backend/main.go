package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"orcaai/backend/config"
	"orcaai/backend/database"
	"orcaai/backend/handlers"
	"orcaai/backend/middleware"
	"orcaai/backend/models"
	"orcaai/backend/orchestrator"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Set production mode if not in debug
	if os.Getenv("DEBUG") != "true" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Load environment variables - Docker'da .env dosyası gerekmeyebilir
	// Önce env dosyası olmadan deneyeceğiz, yoksa sadece uyarı vereceğiz
	envFiles := []string{
		".env",
		"/root/.env",
		"../.env",
		"/app/.env",
	}

	loaded := false
	for _, envFile := range envFiles {
		if err := godotenv.Load(envFile); err == nil {
			log.Printf("✅ Loaded .env file from %s", envFile)
			loaded = true
			break
		}
	}

	if !loaded {
		log.Println("ℹ️  No .env file found, using environment variables from Docker")
	}

	// Initialize configuration
	cfg := config.Load()

	// Log configuration for debugging (Don't log sensitive info in production)
	if os.Getenv("DEBUG") == "true" {
		log.Printf("Database URL: %s", cfg.DatabaseURL)
		log.Printf("Redis URL: %s", cfg.RedisURL)
		log.Printf("Port: %s", cfg.Port)
		log.Printf("Python Worker URL: %s", cfg.PythonWorkerURL)
	}

	// Initialize database connection with retry logic
	var db *sql.DB
	var err error

	log.Println("🔄 Connecting to database...")
	db, err = database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("✅ Database connection established")

	// Initialize Redis connection
	log.Println("🔄 Connecting to Redis...")
	redis := database.ConnectRedis(cfg.RedisURL)
	defer redis.Close()
	log.Println("✅ Redis connection established")

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
			cacheAddr = "redis:6379"
		}
		cachePassword = os.Getenv("CACHE_REDIS_PASSWORD")
		// cacheDB is 0 by default
	}

	if err := orchestrator.InitializeCache(cacheType, cacheAddr, cachePassword, cacheDB); err != nil {
		log.Printf("⚠️  Warning: Failed to initialize cache: %v", err)
	} else {
		log.Printf("✅ Cache initialized successfully: %s", cacheType)
	}
	defer orchestrator.CloseCache()

	// Initialize metrics
	orchestrator.InitializeMetrics()
	log.Println("✅ Metrics initialized")

	// Initialize Gin router
	r := gin.Default()

	// Apply middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.SecurityHeadersMiddleware())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "orcaai-backend",
			"version":   "1.0.0",
			"timestamp": time.Now().Unix(),
		})
	})

	// Ready check endpoint
	r.GET("/ready", func(c *gin.Context) {
		// Check database connection
		if err := database.DB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"error":  "database connection failed",
			})
			return
		}

		// Check Redis connection
		ctx := context.Background()
		if _, err := database.Redis.Ping(ctx).Result(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"error":  "redis connection failed",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
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
		protected.Use(middleware.RateLimitMiddleware(120))
		protected.Use(middleware.IdempotencyMiddleware())
		{
			// AI query endpoint
			protected.POST("/ai/query", handlers.AIQuery)
			protected.POST("/ai/query/stream", handlers.AIQueryStream)
			protected.GET("/ai/providers", handlers.GetProviders)

			// User management (user can view their own profile)
			protected.GET("/user/profile", handlers.GetProfile)
			protected.PUT("/user/profile", handlers.UpdateProfile)

			// API Key management (user can manage their own API keys)
			protected.GET("/keys", handlers.GetAPIKeys)
			protected.POST("/keys", handlers.CreateAPIKey)
			protected.DELETE("/keys/:id", handlers.DeleteAPIKey)
			protected.POST("/keys/:id/rotate", handlers.RotateAPIKey)

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
				admin.POST("/providers/key", handlers.SetProviderKey)
				admin.GET("/providers/key/status", handlers.GetProviderKeyStatus)
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
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("🌐 Health check available at: http://localhost:%s/health", port)
	log.Printf("📊 Metrics available at: http://localhost:%s/metrics", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
