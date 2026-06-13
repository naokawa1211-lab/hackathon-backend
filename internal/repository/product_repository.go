package repository

import (
	"time"
	"hackathon-backend/internal/db"
	"hackathon-backend/internal/model"
)

// SaveProduct は新しい商品情報を MySQL に保存します
func SaveProduct(prod *model.Product) error {
	query := `
		INSERT INTO products (title, description, price, category, image_url_1, seller_id, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	// デフォルト値や初期値を Go 側でもセットしておく
	if prod.Status == "" {
		prod.Status = "available"
	}

	result, err := db.DB.Exec(
		query, 
		prod.Title, 
		prod.Description, 
		prod.Price, 
		prod.Category, 
		prod.ImageURL1, 
		prod.SellerID, 
		prod.Status,
	)
	if err != nil {
		return err
	}

	// 自動採番された ID を構造体にセット
	id, err := result.LastInsertId()
	if err == nil {
		prod.ID = int(id)
	}

	prod.CreatedAt = time.Now()
	return nil
}