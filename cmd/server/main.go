package main

import (
	"fmt"
	"log"

	"hackathon-backend/internal/config"
	"hackathon-backend/internal/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// 環境変数の読み込み
	config.LoadConfig()

	// データベース接続
	db.ConnectDB()
	defer db.DB.Close() // main関数終了時にデータベースを閉じる defer文

	// Webサーバーのルーターを用意
	r := gin.Default()

	// 疎通確認用のテストAPI (http://localhost:8080/ping)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "バックエンドが応答",
		})//gin.HはJSONを作成、func(c *gin.Context){}は無名関数
	})

	// サーバーを指定ポートで起動
	port := config.Env.ServerPort
	log.Println("サーバー起動:", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("サーバーの起動に失敗: %v", err)
	}
}