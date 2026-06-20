package repository

import (
	"errors" // 💡 エラー生成のために追加
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

// BuyProduct は商品の購入処理をトランザクションで実行します
// 💡 インターフェースや構造体を使わず、パッケージ関数として定義
func BuyProduct(productID int, buyerID string) error {
	// ⚡ トランザクションの開始（グローバルの db.DB を使用）
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // エラー時は自動巻き戻し

	// 💡 A. その商品の現在の「売り手」と「ステータス」をチェック
	var productTitle string
	var sellerID string
	var currentStatus string
	err = tx.QueryRow("SELECT title, seller_id, status FROM products WHERE id = ?", productID).Scan(&productTitle, &sellerID, &currentStatus)
	if err != nil {
		return err
	}

	// 二重購入や自作自演のガード
	if currentStatus == "sold" {
		return errors.New("this product is already sold")
	}
	if buyerID == sellerID {
		return errors.New("cannot buy your own product")
	}

	// 💡 B. transactions テーブルに購入履歴をインサート
	_, err = tx.Exec(
		"INSERT INTO transactions (product_id, buyer_id, seller_id, status) VALUES (?, ?, ?, 'completed')",
		productID, buyerID, sellerID,
	)
	if err != nil {
		return err
	}

	// 💡 C. products テーブルのステータスを 'sold' に更新
	_, err = tx.Exec("UPDATE products SET status = 'sold' WHERE id = ?", productID)
	if err != nil {
		return err
	}

	autoContent := "商品【" + productTitle + "】の取引が完了しました。これよりトークルームを解放します。"
	
	_, err = tx.Exec(
		"INSERT INTO messages (sender_id, receiver_id, content) VALUES (?, ?, ?)",
		buyerID, sellerID, autoContent,
	)
	if err != nil {
		return err // メッセージの保存に失敗した場合も、購入ごとロールバックして安全を担保
	}
	// ⚡ すべて成功したらコミット（確定）
	return tx.Commit()
}