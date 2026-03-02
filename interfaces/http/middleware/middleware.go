package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/adrianland/mi-proyecto-blog-api/interfaces/dto"
	"github.com/gin-gonic/gin"
)

// RateLimiter implementa rate limiting por IP
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Limpiar requests antiguos cada minuto
	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, times := range rl.requests {
			filtered := make([]time.Time, 0)
			for _, t := range times {
				if now.Sub(t) < rl.window {
					filtered = append(filtered, t)
				}
			}
			if len(filtered) == 0 {
				delete(rl.requests, ip)
			} else {
				rl.requests[ip] = filtered
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	times, exists := rl.requests[ip]

	if !exists {
		rl.requests[ip] = []time.Time{now}
		return true
	}

	// Limpiar requests antiguos
	filtered := make([]time.Time, 0)
	for _, t := range times {
		if now.Sub(t) < rl.window {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) < rl.limit {
		filtered = append(filtered, now)
		rl.requests[ip] = filtered
		return true
	}

	return false
}

// RateLimitMiddleware middleware para rate limiting
func RateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, dto.ErrorResponse{
				Error:   "TOO_MANY_REQUESTS",
				Message: "Rate limit exceeded. Maximum requests per minute reached.",
				Code:    http.StatusTooManyRequests,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SecurityHeadersMiddleware añade headers de seguridad
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Protección contra ataques
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")

		// Limitar tamaño de request
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20) // 10 MB

		c.Next()
	}
}

// CORSMiddleware configura CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// LoggingMiddleware middleware de logging
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime)
		fmt.Printf("[%s] %s %s %d %v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			c.Request.Method,
			c.Request.RequestURI,
			c.Writer.Status(),
			duration,
		)
	}
}

// RecoveryMiddleware recupera de panics
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error:   "INTERNAL_SERVER_ERROR",
					Message: "An unexpected error occurred",
					Code:    http.StatusInternalServerError,
				})
			}
		}()
		c.Next()
	}
}
