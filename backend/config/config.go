package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
	Port        string

	// AI Providers
	OpenAIAPIKey string
	ClaudeAPIKey string
	GeminiAPIKey string

	// Cache settings
	CacheEnabled    bool
	CacheExpiration time.Duration

	// Rate limiting
	RateLimit       int
	RateLimitWindow time.Duration

	// Monitoring
	PrometheusEnabled bool
	LogLevel          string

	// Python worker
	PythonWorkerURL string
}

func Load() *Config {
	// Debug environment variables
	fmt.Printf("DEBUG: Environment variables:\n")
	fmt.Printf("  DB_HOST: %s\n", os.Getenv("DB_HOST"))
	fmt.Printf("  DB_PORT: %s\n", os.Getenv("DB_PORT"))
	fmt.Printf("  DB_NAME: %s\n", os.Getenv("DB_NAME"))
	fmt.Printf("  DB_USER: %s\n", os.Getenv("DB_USER"))
	fmt.Printf("  REDIS_ADDR: %s\n", os.Getenv("REDIS_ADDR"))
	fmt.Printf("  PORT: %s\n", os.Getenv("PORT"))

	// Construct database URL from environment variables if they exist
	dbHost := getEnv("DB_HOST", "postgres")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "orcaaidb")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "123456789")

	// Construct the database URL from components
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	fmt.Printf("DEBUG: Constructed database URL: %s\n", databaseURL)

	// Redis URL
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisURL := "redis://" + redisAddr

	fmt.Printf("DEBUG: Constructed Redis URL: %s\n", redisURL)

	return &Config{
		DatabaseURL: databaseURL,
		RedisURL:    redisURL,
		JWTSecret:   getEnv("JWT_SECRET", "your-jwt-secret-here"),
		Port:        getEnv("PORT", "8080"),

		// AI Providers
		OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),
		ClaudeAPIKey: getEnv("CLAUDE_API_KEY", ""),
		GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),

		// Cache
		CacheEnabled:    getBoolEnv("CACHE_ENABLED", true),
		CacheExpiration: getDurationEnv("CACHE_EXPIRATION", 24*time.Hour),

		// Rate limiting
		RateLimit:       getIntEnv("RATE_LIMIT", 100),
		RateLimitWindow: getDurationEnv("RATE_LIMIT_WINDOW", time.Hour),

		// Monitoring
		PrometheusEnabled: getBoolEnv("PROMETHEUS_ENABLED", true),
		LogLevel:          getEnv("LOG_LEVEL", "info"),

		// Python worker
		PythonWorkerURL: getEnv("PYTHON_WORKER_URL", "http://ai-worker:8001"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
