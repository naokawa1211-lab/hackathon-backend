package router

import (
    "hackathon-backend/internal/handler"
    "github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    // 💡 2. 手書きの c.Writer.Header()... は全部消して、これに置き換える！
    // ハッカソン用に、すべてのOrigin、すべてのメソッド、すべてのヘッダーを完全に解放する最強設定です
    r.Use(func(c *gin.Context) {
        origin := c.GetHeader("Origin")
        if origin != "" {
            c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
        } else {
            c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        }
        
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        
        // 💡 ★ここの一番最後の部分に「, x-user-id」を追加しました！
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-user-id")
        
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
    
    r.POST("/api/products", handler.CreateProductHandler)
    r.GET("/api/products", handler.GetProductsHandler)
    r.POST("/api/products/:id/buy", handler.BuyProductHandler)
    r.POST("/api/ai/space-description", handler.GenerateSpaceDescriptionHandler)
    r.GET("/api/messages/partners", handler.GetChatPartnersHandler)

    r.POST("/api/users", handler.UpsertUserHandler)     // POST: 同期用
	r.GET("/api/users/:id", handler.GetUserHandler)
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "success", "message": "バックエンドが応答"})
    })

    return r
}