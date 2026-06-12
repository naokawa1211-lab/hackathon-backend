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

var Env *Config //Envをポインタとして定義

// LoadConfig は.envから環境変数を読み込みます
func LoadConfig() {
	// ローカル開発時のみ .env を読み込む（本番環境ではエラーを無視して直接環境変数を取るため）
	if err := godotenv.Load(); err != nil {
		log.Println("Info: .env file not found. Using system environment variables.")
	}

	Env = &Config{
		DBUser:     getEnv("DB_USER", "fallback_user"),
		DBPassword: getEnv("DB_PASSWORD", "fallback_password"),
		DBHost:     getEnv("DB_HOST", "fallback_host"),
		DBName:     getEnv("DB_NAME", "fallback_name"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value // .envにあれば返す
	}
	return fallback //なければデフォルト値を返す
}