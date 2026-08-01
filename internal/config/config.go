package config

import "os"

type Config struct {
	Env          string
	Port         string
	JWTSecret    string
	PostgresURL  string
	RedisURL     string
	UploadDir    string
	AllowedOrigin string
}

func Load() Config {
	return Config{
		Env:          getenv("ENV", "development"),
		Port:         getenv("PORT", "8080"),
		JWTSecret:    getenv("JWT_SECRET", "dev-secret-change-in-production"),
		PostgresURL:  getenv("POSTGRES_URL", "postgres://localhost:5432/shadowchat?sslmode=disable"),
		RedisURL:     getenv("REDIS_URL", "redis://localhost:6379"),
		UploadDir:    getenv("UPLOAD_DIR", "/data/uploads"),
		AllowedOrigin: getenv("ALLOWED_ORIGIN", "*"),
	}
}

func getenv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
