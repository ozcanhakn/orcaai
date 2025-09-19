package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtKey = []byte("orcaai_secret_key") // In production, use environment variable

// Claims represents the JWT claims
type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	IsActive bool      `json:"is_active"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		// Check for parsing errors
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if the token is valid
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if user is active
		if !claims.IsActive {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User account is deactivated"})
			c.Abort()
			return
		}

		// Check if token has expired
		if claims.ExpiresAt.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("is_active", claims.IsActive)

		c.Next()
	}
}

// RateLimitMiddleware provides simple IP-based rate limiting
func RateLimitMiddleware(limitPerMinute int) gin.HandlerFunc {
    type bucket struct{ count int; reset int64 }
    var store = make(map[string]*bucket)
    return func(c *gin.Context) {
        key := c.ClientIP()
        now := time.Now().Unix()
        b, ok := store[key]
        if !ok || now >= b.reset {
            store[key] = &bucket{count: 0, reset: now + 60}
            b = store[key]
        }
        if b.count >= limitPerMinute {
            c.Header("Retry-After", "60")
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
            return
        }
        b.count++
        c.Next()
    }
}

// IdempotencyMiddleware deduplicates POST requests with Idempotency-Key header
func IdempotencyMiddleware() gin.HandlerFunc {
    type entry struct{ status int; body []byte; ts time.Time }
    var cache = make(map[string]*entry)
    return func(c *gin.Context) {
        if c.Request.Method != http.MethodPost {
            c.Next()
            return
        }
        key := c.GetHeader("Idempotency-Key")
        if key == "" {
            c.Next()
            return
        }
        if e, ok := cache[key]; ok && time.Since(e.ts) < 5*time.Minute {
            c.Writer.WriteHeaderNow()
            c.Writer.Write(e.body)
            c.Abort()
            return
        }
        rw := &bufferedWriter{ResponseWriter: c.Writer}
        c.Writer = rw
        c.Next()
        cache[key] = &entry{status: rw.status, body: rw.buf, ts: time.Now()}
    }
}

type bufferedWriter struct {
    gin.ResponseWriter
    buf    []byte
    status int
}

func (w *bufferedWriter) WriteHeader(code int) {
    w.status = code
    w.ResponseWriter.WriteHeader(code)
}

// SecurityHeadersMiddleware adds common security headers
func SecurityHeadersMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        c.Header("X-XSS-Protection", "0")
        c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
        // Basic CSP; consider tightening in production
        c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; img-src 'self' data:; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-eval' 'unsafe-inline'")
        c.Next()
    }
}

func (w *bufferedWriter) Write(b []byte) (int, error) {
    w.buf = append(w.buf, b...)
    return w.ResponseWriter.Write(b)
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(userID uuid.UUID, email, role string, isActive bool) (string, error) {
	// Set expiration time (24 hours)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims
	claims := &Claims{
		UserID:   userID,
		Email:    email,
		Role:     role,
		IsActive: isActive,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID.String(),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// RefreshToken creates a new JWT token with extended expiration
func RefreshToken(userID uuid.UUID, email, role string, isActive bool) (string, error) {
	return GenerateToken(userID, email, role, isActive)
}
