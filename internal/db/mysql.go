package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"hackathon-backend/internal/config"

	_ "github.com/go-sql-driver/mysql" // MySQLドライバーを初期化
)

var DB *sql.DB

// ConnectDB はMySQLデータベースへの接続を確立します
func ConnectDB() {
	// DSN (Data Source Name) を環境変数から組み立て
	// 例: user:password@tcp(127.0.0.1:3306)/dbname
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Asia%%2FTokyo",
		config.Env.DBUser,
		config.Env.DBPassword,
		config.Env.DBHost,
		config.Env.DBPort,
		config.Env.DBName,
	)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ MySQLを開くのに失敗しました: %v", err)
	}

	// ハッカソンや本番（Cloud SQL）で安全に動かすためのコネクション管理設定
	DB.SetMaxOpenConns(20)                 // 同時に開ける最大接続数
	DB.SetMaxIdleConns(5)                  // プール内に保持する空き接続数
	DB.SetConnMaxLifetime(5 * time.Minute) // 接続の寿命

	// 実際に通信できるか生存確認（Ping）
	if err := DB.Ping(); err != nil {
		log.Fatalf("❌ MySQLへの通信（Ping）が届きません。Dockerが起動しているか確認してください: %v", err)
	}

	log.Println("🌌 MySQL データベースへの接続に成功しました！")
}