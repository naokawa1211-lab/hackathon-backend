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

	// Cloudflare R2 (S3互換オブジェクトストレージ) 設定
	R2BucketName      string
	R2AccountID       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2PublicURL       string

	GeminiAPIKey string
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

		R2BucketName:      getEnv("R2_BUCKET_NAME", ""),
		R2AccountID:       getEnv("R2_ACCOUNT_ID", ""),
		R2AccessKeyID:     getEnv("R2_ACCESS_KEY_ID", ""),
		R2SecretAccessKey: getEnv("R2_SECRET_ACCESS_KEY", ""),
		R2PublicURL:       getEnv("R2_PUBLIC_URL", ""),

		GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value // .envにあれば返す
	}
	return fallback //なければデフォルト値を返す
}