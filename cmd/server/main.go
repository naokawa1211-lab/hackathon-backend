package main

import (
	"fmt"
	"log"

	"hackathon-backend/internal/config"
	"hackathon-backend/internal/db"
	"hackathon-backend/internal/router"
)

func main() {
	// 環境変数の読み込み
	config.LoadConfig()

	// データベース接続
	db.ConnectDB()
	defer db.DB.Close() // main関数終了時にデータベースを閉じる defer文

	// Webサーバーのルーター
	r := router.SetupRouter()

	// サーバーを指定ポートで起動
	port := config.Env.ServerPort
	log.Println("サーバー起動:", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("サーバーの起動に失敗: %v", err)
	}
}