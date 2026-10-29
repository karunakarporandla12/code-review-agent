package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                  string
	DatabaseURL           string
	JWTSecret             string
	EncryptionKey         string
	EnableEncryption      bool
	CORSOrigins           []string
	WorkDir               string
	MaxUploadMB           int64
	GitHubWebhookSecret   string
	AzureOpenAIEnabled    bool
	AzureOpenAIEndpoint   string
	AzureOpenAIAPIKey     string
	AzureOpenAIDeployment string
	AzureOpenAIAPIVersion string
}

func Load() Config {
	_ = godotenv.Load()
	maxMB, _ := strconv.ParseInt(env("MAX_UPLOAD_MB", "20"), 10, 64)
	key := env("APP_ENCRYPTION_KEY", env("JWT_SECRET", "change-this-secret"))
	return Config{
		Port:                  env("PORT", "8080"),
		DatabaseURL:           env("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/code_review_agent?sslmode=disable"),
		JWTSecret:             env("JWT_SECRET", "change-this-secret"),
		EncryptionKey:         key,
		EnableEncryption:      env("APP_ENABLE_ENCRYPTION", "true") == "true",
		CORSOrigins:           splitCSV(env("CORS_ORIGINS", defaultCORSOrigins())),
		WorkDir:               env("WORK_DIR", "./data/work"),
		MaxUploadMB:           maxMB,
		GitHubWebhookSecret:   env("GITHUB_WEBHOOK_SECRET", ""),
		AzureOpenAIEnabled:    env("AZURE_OPENAI_ENABLED", "false") == "true",
		AzureOpenAIEndpoint:   strings.TrimRight(env("AZURE_OPENAI_ENDPOINT", ""), "/"),
		AzureOpenAIAPIKey:     env("AZURE_OPENAI_API_KEY", ""),
		AzureOpenAIDeployment: env("AZURE_OPENAI_DEPLOYMENT", ""),
		AzureOpenAIAPIVersion: env("AZURE_OPENAI_API_VERSION", "2024-10-21"),
	}
}

func defaultCORSOrigins() string {
	return "http://localhost:3000,http://localhost:3001,http://localhost:3002,http://localhost:5173,http://127.0.0.1:3000,http://127.0.0.1:3001,http://127.0.0.1:3002,http://127.0.0.1:5173"
}

func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if cleaned := strings.TrimSpace(part); cleaned != "" {
			result = append(result, cleaned)
		}
	}
	return result
}
