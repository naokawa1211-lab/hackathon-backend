package handler

import (
	"net/http"

	"hackathon-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type CreateReviewRequest struct {
	TransactionID int    `json:"transaction_id" binding:"required"`
	ReviewerID    string `json:"reviewer_id" binding:"required"`
	RevieweeID    string `json:"reviewee_id" binding:"required"`
	Rating        int    `json:"rating" binding:"required"`
	Comment       string `json:"comment"`
}

// CreateReviewHandler は購入済みの取引に対するレビュー（評価）を登録します
func CreateReviewHandler(c *gin.Context) {
	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "transaction_id, reviewer_id, reviewee_id, rating は必須です"})
		return
	}

	if req.Rating < 1 || req.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rating は1〜5で指定してください"})
		return
	}

	if err := repository.CreateReview(req.TransactionID, req.ReviewerID, req.RevieweeID, req.Rating, req.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "この取引には既にレビューを投稿済みか、登録に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "レビューを登録しました"})
}

// GetReviewsHandler は指定したユーザーが受け取ったレビューの一覧を返します
func GetReviewsHandler(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id は必須です"})
		return
	}

	reviews, err := repository.GetReviewsForUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "レビューの取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}
