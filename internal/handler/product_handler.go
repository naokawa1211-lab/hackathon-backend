package handler

import (
	"net/http"
	"hackathon-backend/internal/model"
	"hackathon-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

// CreateProductRequest はフロントから届くJSONのバリデーション定義です
type CreateProductRequest struct {
	Title         string `json:"title" binding:"required"`
	Description   string `json:"description"`
	Price         int    `json:"price" binding:"required"`
	Category      string `json:"category" binding:"required"`
	ImageURL1     string `json:"image_url_1"`
	SellerID      string `json:"seller_id"` // ハッカソン中は一旦モックIDでもOK
}

func CreateProductHandler(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 初期画像は銀河
	if req.ImageURL1 == "" {
		req.ImageURL1 = "https://images.unsplash.com/photo-1506318137071-a8e063b4bec0?q=80&w=600" // カッコいい銀河のフリー画像
	}

	product := model.Product{
		Title:         req.Title,
		Description:   req.Description,
		Price:         req.Price,
		Category:      req.Category,
		ImageURL1: 	   req.ImageURL1,
		SellerID:      req.SellerID,
	}

	// repository 層に保存を丸投げ
	if err := repository.SaveProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "商品情報の登録に失敗しました"})
		return
	}

	// 成功したら 200 OK と登録されたデータを返す
	c.JSON(http.StatusOK, product)
}