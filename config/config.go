package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DBType        string
	DBUrl         string
	JWTSecret     string
	SendGridKey   string
	EmailFrom     string
	WeatherAPIKey string
	BaseURL       string
}

var C *Config

func LoadConfig() {
	_ = godotenv.Load()

	C = &Config{
		Port:          getEnv("PORT", "8080"),
		DBType:        getEnv("DB_TYPE", "postgres"),
		DBUrl:         mustGet("DB_URL"),
		JWTSecret:     mustGet("JWT_SECRET"),
		SendGridKey:   mustGet("SENDGRID_API_KEY"),
		EmailFrom:     getEnv("EMAIL_FROM", "weatherapp-no-reply@woolberry.ua"),
		WeatherAPIKey: mustGet("WEATHER_API_KEY"),
		BaseURL:       strings.TrimRight(getEnv("BASE_URL", "http://localhost:8080"), "/"),
	}
}

func Reload() {
	C = nil
	LoadConfig()
}

func mustGet(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return val
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
