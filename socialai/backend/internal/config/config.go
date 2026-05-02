package config

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Environment       string
	Address           string
	PublicURL         string
	AllowedOrigins    []string
	TrustProxy        bool
	JWTSecret         []byte
	TokenTTL          time.Duration
	RepositoryBackend string
	DataFile          string
	ElasticsearchURL  string
	ElasticsearchUser string
	ElasticsearchPass string
	ElasticsearchPref string
	MediaBackend      string
	MediaDirectory    string
	MaxUploadBytes    int64
	GCSBucket         string
	GCSBearerToken    string
	AIBackend         string
	OpenAIAPIKey      string
	OpenAIModel       string
	OpenAIBaseURL     string
}

func Load() (Config, error) {
	cfg := Config{
		Environment:       strings.ToLower(env("APP_ENV", "development")),
		Address:           ":" + env("PORT", "8080"),
		PublicURL:         strings.TrimRight(env("PUBLIC_URL", "http://localhost:8080"), "/"),
		AllowedOrigins:    split(env("CORS_ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:3000")),
		TrustProxy:        boolValue("TRUST_PROXY", false),
		TokenTTL:          duration("JWT_TTL", 24*time.Hour),
		RepositoryBackend: env("REPOSITORY_BACKEND", "file"),
		DataFile:          env("DATA_FILE", "../data/socialai.json"),
		ElasticsearchURL:  strings.TrimRight(os.Getenv("ELASTICSEARCH_URL"), "/"),
		ElasticsearchUser: os.Getenv("ELASTICSEARCH_USERNAME"),
		ElasticsearchPass: os.Getenv("ELASTICSEARCH_PASSWORD"),
		ElasticsearchPref: env("ELASTICSEARCH_INDEX_PREFIX", "socialai"),
		MediaBackend:      env("MEDIA_BACKEND", "local"),
		MediaDirectory:    env("MEDIA_DIRECTORY", "../media"),
		MaxUploadBytes:    int64Value("MAX_UPLOAD_BYTES", 25<<20),
		GCSBucket:         os.Getenv("GCS_BUCKET"),
		GCSBearerToken:    os.Getenv("GCS_BEARER_TOKEN"),
		AIBackend:         env("AI_IMAGE_BACKEND", "local"),
		OpenAIAPIKey:      os.Getenv("OPENAI_API_KEY"),
		OpenAIModel:       env("OPENAI_IMAGE_MODEL", "gpt-image-2"),
		OpenAIBaseURL:     strings.TrimRight(env("OPENAI_BASE_URL", "https://api.openai.com/v1"), "/"),
	}
	if cfg.Environment != "development" && cfg.Environment != "test" && cfg.Environment != "production" {
		return Config{}, fmt.Errorf("APP_ENV must be development, test, or production")
	}
	secret := os.Getenv("JWT_SECRET")
	isPlaceholder := secret == "replace-with-at-least-32-random-characters" || secret == "change-me" || secret == "your-secret-here"
	if len(secret) < 32 || isPlaceholder {
		if cfg.Environment == "production" {
			return Config{}, fmt.Errorf("JWT_SECRET must be a non-placeholder secret containing at least 32 characters in production")
		}
		random := make([]byte, 32)
		if _, err := rand.Read(random); err != nil {
			return Config{}, fmt.Errorf("generate development JWT secret: %w", err)
		}
		secret = base64.RawURLEncoding.EncodeToString(random)
	}
	cfg.JWTSecret = []byte(secret)
	if cfg.MaxUploadBytes < 1<<20 {
		return Config{}, fmt.Errorf("MAX_UPLOAD_BYTES must be at least 1048576")
	}
	return cfg, nil
}

func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func split(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}

func duration(key string, fallback time.Duration) time.Duration {
	if parsed, err := time.ParseDuration(os.Getenv(key)); err == nil && parsed > 0 {
		return parsed
	}
	return fallback
}

func int64Value(key string, fallback int64) int64 {
	parsed, err := strconv.ParseInt(os.Getenv(key), 10, 64)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func boolValue(key string, fallback bool) bool {
	parsed, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return parsed
}
