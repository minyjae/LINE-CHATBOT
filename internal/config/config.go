package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                        string
	AppPort                       string
	AppURL                        string
	AppPrefix                     string
	CorsAllows                    string
	DBHost                        string
	DBPort                        string
	DBUser                        string
	DBPassword                    string
	DBName                        string
	DBSSLMode                     string
	JWTSecret                     string
	JWTExpiresIn                  string
	LineChannelSecret             string
	LineChannelAccessToken        string
	OpenAIAPIKey                  string
	OpenAIIntentModel             string
	ReminderWorkerEnabled         bool
	ReminderWorkerIntervalSeconds int
	ReminderWorkerBatchSize       int
	AdminEmail                    string
	AdminFirstName                string
	AdminLastName                 string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variable")
	}

	config := &Config{
		// ค่าปลอดภัย
		AppEnv:       getEnv("APP_ENV", "development"),
		AppPort:      getEnv("APP_PORT", "8080"),
		AppURL:       getEnv("APP_URL", "http://localhost:8080"),
		AppPrefix:    getEnv("APP_PREFIX", ""),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5431"),
		DBUser:       getEnv("DB_USER", "test"),
		DBSSLMode:    getEnv("DB_SSL", "disable"),
		JWTExpiresIn: getEnv("JWT_EXPIRES_IN", "24h"),
		CorsAllows:   getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),

		// ค่าไม่ปลอดภัย
		DBPassword:                    getEnv("DB_PASSWORD", ""),
		DBName:                        getEnv("DB_NAME", ""),
		JWTSecret:                     getEnv("JWT_SECRET", ""),
		LineChannelSecret:             getEnv("LINE_CHANNEL_SECRET", ""),
		LineChannelAccessToken:        getEnv("LINE_CHANNEL_ACCESS_TOKEN", ""),
		OpenAIAPIKey:                  getEnv("OPENAI_API_KEY", ""),
		OpenAIIntentModel:             getEnv("OPENAI_INTENT_MODEL", "gpt-5"),
		ReminderWorkerEnabled:         getEnvBool("REMINDER_WORKER_ENABLED", false),
		ReminderWorkerIntervalSeconds: getEnvInt("REMINDER_WORKER_INTERVAL_SECONDS", 60),
		ReminderWorkerBatchSize:       getEnvInt("REMINDER_WORKER_BATCH_SIZE", 50),
		AdminEmail:                    getEnv("ADMIN_EMAIL", ""),
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(config *Config) error {
	if config.AppEnv == "production" {
		if config.DBPassword == "" {
			return fmt.Errorf("DB_PASS is required for production environment")
		}
		if config.JWTSecret == "" {
			return fmt.Errorf("JWT_SECRET is required for production environment")
		}
		if len(config.JWTSecret) < 32 {
			return fmt.Errorf("JWT_SECRET must be at least 32 charecters long for production")
		}
		if config.DBSSLMode == "disable" {
			log.Println("Warnig: SSL is disabled for database connection in production")
		}
		if config.AdminEmail == "" {
			return fmt.Errorf("ADMIN_EMAIL is required for production environment")
		}
		if config.LineChannelSecret == "" {
			return fmt.Errorf("LINE_CHANNEL_SECRET is required for production environment")
		}
		if config.LineChannelAccessToken == "" {
			return fmt.Errorf("LINE_CHANNEL_ACCESS_TOKEN is required for production environment")
		}
	}

	if config.DBName == "" {
		return fmt.Errorf("DB_Name is required for production environment")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}
