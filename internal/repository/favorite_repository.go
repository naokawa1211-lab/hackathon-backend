package repository

import (
	"database/sql"

	"hackathon-backend/internal/db"
	"hackathon-backend/internal/model"
)

// ToggleFavorite はお気に入りの追加・解除を切り替える。戻り値は切り替え後の状態（true=お気に入り済み）
func ToggleFavorite(userID string, productID int) (bool, error) {
	var existingID int
	err := db.DB.QueryRow("SELECT id FROM favorites WHERE user_id = ? AND product_id = ?", userID, productID).Scan(&existingID)

	if err == sql.ErrNoRows {
		// 未登録なので追加する
		if _, err := db.DB.Exec("INSERT INTO favorites (user_id, product_id) VALUES (?, ?)", userID, productID); err != nil {
			return false, err
		}
		return true, nil
	}
	if err != nil {
		return false, err
	}

	// 既に登録済みなので解除する
	if _, err := db.DB.Exec("DELETE FROM favorites WHERE id = ?", existingID); err != nil {
		return false, err
	}
	return false, nil
}

// GetFavoriteProducts は指定したユーザーがお気に入り登録した商品の一覧を取得する
func GetFavoriteProducts(userID string) ([]model.Product, error) {
	query := `
		SELECT p.id, p.title, p.description, p.price, p.category, p.image_url_1, p.seller_id, p.status
		FROM favorites f
		JOIN products p ON f.product_id = p.id
		WHERE f.user_id = ?
		ORDER BY f.created_at DESC
	`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Price, &p.Category, &p.ImageURL1, &p.SellerID, &p.Status); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
