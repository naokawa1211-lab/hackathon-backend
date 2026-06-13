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

// GetAllProducts - DBから商品一覧を全件取得する
func GetAllProducts() ([]model.Product, error) {
	var products []model.Product

	// ここで repository が使っている実際のDBインスタンス（例: db.DB や r.db など）から取得します。
	// すでに同じファイルにある SaveProduct の実装を参考に、クエリ部分を以下のように合わせます。
	// 以下は sqlx や database/sql を想定した標準的な例です：
	rows, err := db.DB.Query("SELECT id, title, description, price, category, image_url_1, seller_id, status FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Price, &p.Category, &p.ImageURL1, &p.SellerID, &p.Status); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}