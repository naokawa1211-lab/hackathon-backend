package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"hackathon-backend/internal/config"

	_ "github.com/go-sql-driver/mysql" // MySQLドライバーを裏で初期化する _をつけないとエラーになる
)

var DB *sql.DB

// ConnectDB はMySQLデータベースへの接続を確立
func ConnectDB() {
	// DSN (Data Source Name) を環境変数から組み立て
	dsn := fmt.Sprintf("%s:%s@%s/%s?parseTime=true",
		config.Env.DBUser, 
		config.Env.DBPassword, 
		config.Env.DBHost, 
		config.Env.DBName,
	)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("MySQLを開くのに失敗しました: %v", err)
	}

	//コネクションプール
	DB.SetMaxOpenConns(20)                 // 同時に開ける最大接続数
	DB.SetMaxIdleConns(5)                  // プール内に保持する空き接続数
	DB.SetConnMaxLifetime(5 * time.Minute) // 接続の寿命

	// Pingで実際に通信を行う
	if err := DB.Ping(); err != nil {
		log.Fatalf("MySQLへの通信（Ping）が届きません。Dockerが起動しているか確認してください: %v", err)
	}

	log.Println("MySQL データベースへの接続に成功しました")
}