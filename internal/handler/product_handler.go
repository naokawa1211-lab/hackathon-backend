package handler

import (
	"log"
	"net/http"
	"strconv" 
	"hackathon-backend/internal/model"
	"hackathon-backend/internal/repository"
	"hackathon-backend/internal/storage"

	"github.com/gin-gonic/gin"
)

// 画像未指定時のフォールバック画像
const defaultProductImageURL = "https://images.unsplash.com/photo-1506318137071-a8e063b4bec0?q=80&w=600"

// CreateProductHandler は multipart/form-data で商品情報と画像を受け取り、
// 画像をR2にアップロードした上で商品をDBに登録します
func CreateProductHandler(c *gin.Context) {
	title := c.PostForm("title")
	category := c.PostForm("category")
	priceStr := c.PostForm("price")

	if title == "" || category == "" || priceStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title, price, category は必須です"})
		return
	}

	price, err := strconv.Atoi(priceStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price は数値で指定してください"})
		return
	}

	description := c.PostForm("description")
	sellerID := c.PostForm("seller_id") // ハッカソン中は一旦モックIDでもOK

	imageURL := defaultProductImageURL
	if fileHeader, ferr := c.FormFile("image"); ferr == nil && fileHeader != nil {
		uploadedURL, uploadErr := storage.UploadImage(fileHeader)
		if uploadErr != nil {
			log.Printf("[CreateProductHandler] R2 upload failed: filename=%q err=%v", fileHeader.Filename, uploadErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "画像のアップロードに失敗しました"})
			return
		}
		imageURL = uploadedURL
	}

	product := model.Product{
		Title:       title,
		Description: description,
		Price:       price,
		Category:    category,
		ImageURL1:   imageURL,
		SellerID:    sellerID,
	}

	if err := repository.SaveProduct(&product); err != nil {
		log.Printf("[CreateProductHandler] SaveProduct failed: title=%q err=%v", product.Title, err)
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
		// 💡 本番で500が起きた際に「接続上限なのか／別の理由か」を後から特定できるよう詳細を残す
		log.Printf("[GetProductsHandler] GetAllProducts failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "商品一覧の取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetPurchasedProductsHandler は指定したユーザーが購入した商品の一覧を返します
func GetPurchasedProductsHandler(c *gin.Context) {
	buyerID := c.Query("buyer_id")
	if buyerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "buyer_id は必須です"})
		return
	}

	purchased, err := repository.GetPurchasedProducts(buyerID)
	if err != nil {
		log.Printf("[GetPurchasedProductsHandler] GetPurchasedProducts failed: buyer_id=%q err=%v", buyerID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "購入履歴の取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, purchased)
}

// BuyProductHandler は商品の購入処理を呼び出します
// 💡 Ginの仕様に合わせて新規追加！
func BuyProductHandler(c *gin.Context) {
	// 1. フロントから送られてきたヘッダー（買い手のUID）をキャッチ
	buyerID := c.GetHeader("X-User-id")
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
		if err.Error() == "this product is already sold" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "この商品は既に売り切れです"})
			return
		}
		if err.Error() == "cannot buy your own product" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "自分の出品した商品は購入できません"})
			return
		}
		// 🚀 【追加】残高不足エラーをキャッチしてフロントに400で伝える！
		if err.Error() == "insufficient credits" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "クレジットが不足しています"})
			return
		}

		log.Printf("[BuyProductHandler] BuyProduct failed: product_id=%d buyer_id=%q err=%v", productID, buyerID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "購入処理に失敗しました"})
		return
	}

	// 4. 成功したら200 OKとメッセージを返す
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "TRANSACTION COMPLETED: 取引が正常に成立しました",
	})
}