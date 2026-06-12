package main

import (
	"fmt"
	"log"

	"hackathon-backend/internal/config"
	"hackathon-backend/internal/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 設定（環境変数）の読み込み
	config.LoadConfig()

	// 2. データベース接続
	db.ConnectDB()
	defer db.DB.Close() // サーバー終了時に安全にDBを閉じる

	// 3. Gin（Webフレームワーク）の起動
	r := gin.Default()

	// 疎通確認用のテストAPI (http://localhost:8080/ping)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "🪐 バックエンド宇宙ステーションからの応答を受信！",
		})
	})

	// 4. サーバーを指定ポートで起動
	port := config.Env.ServerPort
	log.Printf("🚀 バックエンドサーバーがポート %s で離陸しました！", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("❌ サーバーの起動に失敗しました: %v", err)
	}
}