package router

import (
	"hackathon-backend/internal/handler"
	"github.com/gin-gonic/gin"
)

// SetupRouter - ルーティングの初期化
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORSミドルウェアの設定
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API エンドポイント
	r.POST("/api/messages", handler.SendMessageHandler)
	r.GET("/api/messages/history", handler.GetChatHistoryHandler)
	
	// 商品系エンドポイント（GETを追加！）
	r.POST("/api/products", handler.CreateProductHandler)
	r.GET("/api/products", handler.GetProductsHandler)

	// 疎通確認用
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "バックエンドが応答",
		})
	})

	return r
}