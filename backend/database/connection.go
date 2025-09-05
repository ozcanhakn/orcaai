package database

import (
	"database/sql"
	"log"
	"context"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var (
	DB    *sql.DB
	Redis *redis.Client
)

func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	DB = db
	log.Println("✅ Connected to PostgreSQL")

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Printf("⚠️  Migration error: %v", err)
	}

	return db, nil
}

func ConnectRedis(redisURL string) *redis.Client {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	client := redis.NewClient(opt)

	// Test connection
	ctx := context.Background()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	Redis = client
	log.Println("✅ Connected to Redis")
	return client
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
		
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			role VARCHAR(50) DEFAULT 'user',
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		
		`CREATE TABLE IF NOT EXISTS api_keys (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			key_hash VARCHAR(255) UNIQUE NOT NULL,
			last_used_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT NOW(),
			is_active BOOLEAN DEFAULT TRUE
		);`,
		
		`CREATE TABLE IF NOT EXISTS request_logs (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			api_key_id UUID REFERENCES api_keys(id) ON DELETE CASCADE,
			provider VARCHAR(100) NOT NULL,
			model VARCHAR(100),
			prompt_tokens INTEGER DEFAULT 0,
			completion_tokens INTEGER DEFAULT 0,
			cost_usd DECIMAL(10, 6) DEFAULT 0,
			latency_ms INTEGER DEFAULT 0,
			cache_hit BOOLEAN DEFAULT FALSE,
			status VARCHAR(50) DEFAULT 'success',
			error_message TEXT,
			created_at TIMESTAMP DEFAULT NOW()
		);`,
		
		`CREATE TABLE IF NOT EXISTS ai_providers (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(100) UNIQUE NOT NULL,
			base_url VARCHAR(255) NOT NULL,
			api_key_encrypted TEXT,
			cost_per_1k_input DECIMAL(10, 6) DEFAULT 0,
			cost_per_1k_output DECIMAL(10, 6) DEFAULT 0,
			max_tokens INTEGER DEFAULT 4000,
			is_active BOOLEAN DEFAULT TRUE,
			priority INTEGER DEFAULT 1,
			created_at TIMESTAMP DEFAULT NOW()
		);`,
		
		`CREATE TABLE IF NOT EXISTS user_preferences (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			preferred_provider VARCHAR(100),
			max_cost_per_request DECIMAL(10, 6) DEFAULT 1.0,
			enable_caching BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		
		// Indexes for performance
		`CREATE INDEX IF NOT EXISTS idx_request_logs_user_id ON request_logs(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_request_logs_created_at ON request_logs(created_at);`,
		`CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_api_keys_key_hash ON api_keys(key_hash);`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return err
		}
	}

	log.Println("✅ Database migrations completed")
	return nil
}