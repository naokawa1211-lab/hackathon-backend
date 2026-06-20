package handler

import (
	"net/http"

	"hackathon-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type ToggleFavoriteRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	ProductID int    `json:"product_id" binding:"required"`
}

// ToggleFavoriteHandler はお気に入りの追加・解除を切り替えます
func ToggleFavoriteHandler(c *gin.Context) {
	var req ToggleFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id と product_id は必須です"})
		return
	}

	favorited, err := repository.ToggleFavorite(req.UserID, req.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "お気に入りの更新に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"favorited": favorited})
}

// GetFavoritesHandler は指定したユーザーのお気に入り商品一覧を返します
func GetFavoritesHandler(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id は必須です"})
		return
	}

	products, err := repository.GetFavoriteProducts(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "お気に入りの取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, products)
}
