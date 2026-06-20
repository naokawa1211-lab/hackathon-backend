// internal/handler/user_handler.go
package handler

import (
	"net/http"
	"hackathon-backend/internal/model"
	"hackathon-backend/internal/repository"
	"github.com/gin-gonic/gin"
)

// 💡 パッケージレベルでリポジトリを保持（router.goの直接呼び出しスタイルに完全適合）
var userRepo *repository.UserRepository

// SetUserRepository は main.go などからDB接続後にリポジトリを注入するための関数です
func SetUserRepository(repo *repository.UserRepository) {
	userRepo = repo
}
// UpsertUserHandler: フロントから送られてきたユーザー情報をMySQLに保存
func UpsertUserHandler(c *gin.Context) {
	var user model.User
	// Ginの機能でJSONを構造体にバインド
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if user.ID == "" || user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields (id, username)"})
		return
	}

	if err := userRepo.UpsertUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー情報のMySQL同期に失敗しました"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "User synced successfully"})
}

// GetUserHandler: URLパラメータからIDを取得してユーザー情報を返す
func GetUserHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID"})
		return
	}

	user, err := userRepo.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースからのユーザー取得に失敗しました"})
		return
	}

	// DBに本物のユーザーが存在した場合は、そのデータを完璧に返却！
	if user != nil {
		c.JSON(http.StatusOK, user)
		return
	}

	// 🛸 万が一DBに見つからなかった場合のみ、404エラーを返す
	// これにより、フロント（DM画面）側で自動的に「未知の生命体」などのフォールバック表示に流れます
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found in database"})
}