package main

import (
	"fmt"
	"log"

	"hackathon-backend/internal/config"
	"hackathon-backend/internal/db"
	"hackathon-backend/internal/handler"

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

	r.Use(func(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // すべてのオリジンからのアクセスを許可
    c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(204)
        return
    }

    c.Next()
})
	r.POST("/api/messages", handler.SendMessageHandler)      // メッセージ送信
	r.GET("/api/messages/history", handler.GetChatHistoryHandler) // チャット履歴取得
	r.POST("/api/products", handler.CreateProductHandler)//商品情報登録

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