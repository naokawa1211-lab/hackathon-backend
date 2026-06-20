// internal/handler/user_handler.go
package handler

import (
	"encoding/json"
	"net/http"
	"hackathon-backend/internal/model"
	"hackathon-backend/internal/repository"
)

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
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

	// 💡 TODO: repository層の呼び出し (例: err := userRepo.UpsertUser(&user))
	// 現状、エラーがなければ成功として返す
	
	c.JSON(http.StatusOK, gin.H{"message": "User synced successfully"})
}

// GetUserHandler: URLパラメータからIDを取得してユーザー情報を返す
func GetUserHandler(c *gin.Context) {
	// 💡 Ginでは c.Param("id") で "/api/users/:id" のコロンの部分を直接引っこ抜けます！
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID"})
		return
	}

	// 💡 TODO: repository層からデータを取得 (例: user, err := userRepo.GetUserByID(id))
	// ここではDBから取ってきたと仮定したダミーの返却例を書いておきます
	
	// 見つからなかった場合のモック例（バックエンドが未完成の間の生存戦略）
	if id == "mock_uid_naoya" {
		c.JSON(http.StatusOK, gin.H{
			"id": "mock_uid_naoya",
			"username": "Naoya",
		})
		return
	}

	// 本来はDBから取得したデータを返す
	// c.JSON(http.StatusOK, user)
	
	// まだDBにいない本物ユーザーへの一時的なダミー返却
	c.JSON(http.StatusOK, gin.H{
		"id": id,
		"username": "宇宙の迷い人",
	})
}