package main

import (
	"fmt"
	"log"

	"hackathon-backend/internal/config"
	"hackathon-backend/internal/db"
	"hackathon-backend/internal/router"
	"hackathon-backend/internal/storage"
	"hackathon-backend/internal/handler"
	"hackathon-backend/internal/repository"
)

func main() {
	// 環境変数の読み込み
	config.LoadConfig()

	// データベース接続
	db.ConnectDB()
	defer db.DB.Close() // main関数終了時にデータベースを閉じる defer文

	userRepo := repository.NewUserRepository(db.DB)
	handler.SetUserRepository(userRepo)
	
	// Cloudflare R2 (画像アップロード用) クライアントの初期化
	if err := storage.InitR2(); err != nil {
		log.Fatalf("R2クライアントの初期化に失敗: %v", err)
	}

	// Webサーバーのルーター
	r := router.SetupRouter()

	// サーバーを指定ポートで起動
	port := config.Env.ServerPort
	log.Println("サーバー起動:", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("サーバーの起動に失敗: %v", err)
	}
}