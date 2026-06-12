package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	ServerPort string
}

var Env *Config

// LoadConfig は.envから環境変数を読み込みます
func LoadConfig() {
	// ローカル開発時のみ .env を読み込む（本番環境ではエラーを無視して直接環境変数を取るため）
	if err := godotenv.Load(); err != nil {
		log.Println("Info: .env file not found. Using system environment variables.")
	}

	Env = &Config{
		DBUser:     getEnv("DB_USER", "space_operator"),
		DBPassword: getEnv("DB_PASSWORD", "space_password123"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "milkyway_flea_market"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}