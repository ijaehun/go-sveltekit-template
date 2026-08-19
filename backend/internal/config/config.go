// Package config: 환경변수 기반 설정. .env 파일은 docker-compose가 주입하고,
// 로컬 실행 시엔 셸에서 export 하거나 그냥 기본값으로 뜬다.
package config

import (
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	Port           string
	Environment    string   // development | production | ci
	AllowedOrigins []string // CORS 허용 origin (dev에서 vite가 다른 포트일 때 필요)
}

func Load() Config {
	return Config{
		Port:           getEnv("PORT", "8080"),
		Environment:    getEnv("ENVIRONMENT", "development"),
		AllowedOrigins: splitCSV(getEnv("ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:3000")),
	}
}

func (c Config) LogLevel() slog.Level {
	if c.Environment == "production" {
		return slog.LevelInfo
	}
	return slog.LevelDebug
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
