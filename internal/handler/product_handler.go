package handler

import (
	"net/http"
	"strconv" // 💡 追加：文字列のIDを数値に変えるため
	"hackathon-backend/internal/model"
	"hackathon-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

// CreateProductRequest はフロントから届くJSONのバリデーション定義です
type CreateProductRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Price       int    `json:"price" binding:"required"`
	Category    string `json:"category" binding:"required"`
	ImageURL1   string `json:"image_url_1"`
	SellerID    string `json:"seller_id"` // ハッカソン中は一旦モックIDでもOK
}

func CreateProductHandler(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 初期画像は銀河
	if req.ImageURL1 == "" {
		req.ImageURL1 = "https://images.unsplash.com/photo-1506318137071-a8e063b4bec0?q=80&w=600"
	}

	product := model.Product{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		ImageURL1:   req.ImageURL1,
		SellerID:    req.SellerID,
	}

	if err := repository.SaveProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "商品情報の登録に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func GetProductsHandler(c *gin.Context) {
	var products []model.Product

	var err error
	products, err = repository.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "商品一覧の取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// BuyProductHandler は商品の購入処理を呼び出します
// 💡 Ginの仕様に合わせて新規追加！
func BuyProductHandler(c *gin.Context) {
	// 1. フロントから送られてきたヘッダー（買い手のUID）をキャッチ
	buyerID := c.GetHeader("X-User-UID")
	if buyerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: 宇宙市民UIDが未検出です"})
		return
	}

	// 2. URLパラメータから商品ID（ :id ）を取得
	idStr := c.Param("id")
	productID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効な商品IDです"})
		return
	}

	// 3. repository層のトランザクション関数を実行
	if err := repository.BuyProduct(productID, buyerID); err != nil {
		// リポジトリ側のガードに引っかかった場合のエラーハンドリング
		if err.Error() == "this product is already sold" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "この物資はすでに他のセクターの商人が購入済みです"})
			return
		}
		if err.Error() == "cannot buy your own product" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "自作自演の取引は銀河法により禁止されています"})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取引の同期中にエラーが発生しました"})
		return
	}

	// 4. 成功したら200 OKとメッセージを返す
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "TRANSACTION COMPLETED: 取引が正常に成立しました",
	})
}