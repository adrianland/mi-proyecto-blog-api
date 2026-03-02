package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort              string
	ServerEnv               string
	DBHost                  string
	DBPort                  string
	DBUser                  string
	DBPassword              string
	DBName                  string
	DBMaxConnections        int
	DBMaxIdleConnections    int
	DBConnectionMaxLifetime time.Duration
	RateLimitRequests       int
	RateLimitWindow         time.Duration
	DefaultPageSize         int
	MaxPageSize             int
}

func LoadConfig() *Config {
	// Cargar variables de entorno desde .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	cfg := &Config{
		ServerPort:           getEnv("SERVER_PORT", "8080"),
		ServerEnv:            getEnv("SERVER_ENV", "development"),
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "3306"),
		DBUser:               getEnv("DB_USER", "blog_user"),
		DBPassword:           getEnv("DB_PASSWORD", "secure_password"),
		DBName:               getEnv("DB_NAME", "blog_db"),
		DBMaxConnections:     getEnvInt("DB_MAX_CONNECTIONS", 25),
		DBMaxIdleConnections: getEnvInt("DB_MAX_IDLE_CONNECTIONS", 5),
		RateLimitRequests:    getEnvInt("RATE_LIMIT_REQUESTS", 100),
		DefaultPageSize:      getEnvInt("DEFAULT_PAGE_SIZE", 10),
		MaxPageSize:          getEnvInt("MAX_PAGE_SIZE", 100),
	}

	// Parsear duraciones
	lifetime := getEnv("DB_CONNECTION_MAX_LIFETIME", "10m")
	cfg.DBConnectionMaxLifetime, err = time.ParseDuration(lifetime)
	if err != nil {
		cfg.DBConnectionMaxLifetime = 10 * time.Minute
	}

	window := getEnv("RATE_LIMIT_WINDOW", "1m")
	cfg.RateLimitWindow, err = time.ParseDuration(window)
	if err != nil {
		cfg.RateLimitWindow = 1 * time.Minute
	}

	return cfg
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultVal
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid integer value for %s, using default: %d", key, defaultVal)
		return defaultVal
	}
	return value
}
