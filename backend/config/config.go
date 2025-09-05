package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL    string
	RedisURL       string
	JWTSecret      string
	Port           string
	
	// AI Providers
	OpenAIAPIKey     string
	ClaudeAPIKey     string
	GeminiAPIKey     string
	
	// Cache settings
	CacheEnabled     bool
	CacheExpiration  time.Duration
	
	// Rate limiting
	RateLimit        int
	RateLimitWindow  time.Duration
	
	// Monitoring
	PrometheusEnabled bool
	LogLevel         string
	
	// Python worker
	PythonWorkerURL  string
}

func Load() *Config {
	return &Config{
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://localhost/orcaai?sslmode=disable"),
		RedisURL:         getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:        getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this"),
		Port:            getEnv("PORT", "8080"),
		
		// AI Providers
		OpenAIAPIKey:    getEnv("OPENAI_API_KEY", ""),
		ClaudeAPIKey:    getEnv("CLAUDE_API_KEY", ""),
		GeminiAPIKey:    getEnv("GEMINI_API_KEY", ""),
		
		// Cache
		CacheEnabled:    getBoolEnv("CACHE_ENABLED", true),
		CacheExpiration: getDurationEnv("CACHE_EXPIRATION", 24*time.Hour),
		
		// Rate limiting
		RateLimit:       getIntEnv("RATE_LIMIT", 100),
		RateLimitWindow: getDurationEnv("RATE_LIMIT_WINDOW", time.Hour),
		
		// Monitoring
		PrometheusEnabled: getBoolEnv("PROMETHEUS_ENABLED", true),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		
		// Python worker
		PythonWorkerURL: getEnv("PYTHON_WORKER_URL", "http://localhost:8001"),
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