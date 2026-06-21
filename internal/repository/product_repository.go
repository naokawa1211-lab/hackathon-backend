package repository

import (
	"database/sql"
	"errors" // 💡 エラー生成のために追加
	"log"
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
		// 💡 description / category / seller_id はDBスキーマ上NULLを許容するカラム（NOT NULL制約なし）。
		// 直接 string へScanすると、本番データにNULLが1件でも混ざった瞬間に
		// 「sql: Scan error on column index N」で全件取得が失敗する。
		// sql.NullString で安全に受け取り、NULLの場合は空文字にフォールバックする。
		var description, category, sellerID sql.NullString
		if err := rows.Scan(&p.ID, &p.Title, &description, &p.Price, &category, &p.ImageURL1, &sellerID, &p.Status); err != nil {
			// 🔭 診断用: 次回も同じエラーが起きた場合に原因のカラム・型をすぐ特定できるようにする
			if colTypes, ctErr := rows.ColumnTypes(); ctErr == nil {
				for i, ct := range colTypes {
					log.Printf("[GetAllProducts] column %d: name=%s dbType=%s", i, ct.Name(), ct.DatabaseTypeName())
				}
			}
			return nil, err
		}
		p.Description = description.String
		p.Category = category.String
		p.SellerID = sellerID.String
		products = append(products, p)
	}

	return products, nil
}

// GetPurchasedProducts は指定したbuyerIDが購入した商品の一覧を取引情報付きで取得する
// （その取引に対して既にレビュー済みかどうかも合わせて返す）
func GetPurchasedProducts(buyerID string) ([]model.PurchasedProduct, error) {
	query := `
		SELECT t.id, t.created_at, p.id, p.title, p.description, p.price, p.category, p.image_url_1, p.seller_id, p.status, p.created_at,
			EXISTS(SELECT 1 FROM reviews r WHERE r.transaction_id = t.id AND r.reviewer_id = t.buyer_id) AS already_reviewed
		FROM transactions t
		JOIN products p ON t.product_id = p.id
		WHERE t.buyer_id = ?
		ORDER BY t.created_at DESC
	`

	rows, err := db.DB.Query(query, buyerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchased []model.PurchasedProduct
	for rows.Next() {
		var pp model.PurchasedProduct
		// 💡 GetAllProductsと同じ理由でNULL許容カラムを安全に受け取る
		var description, category, sellerID sql.NullString
		if err := rows.Scan(
			&pp.TransactionID, &pp.PurchasedAt,
			&pp.ID, &pp.Title, &description, &pp.Price, &category, &pp.ImageURL1, &sellerID, &pp.Status, &pp.CreatedAt,
			&pp.AlreadyReviewed,
		); err != nil {
			return nil, err
		}
		pp.Description = description.String
		pp.Category = category.String
		pp.SellerID = sellerID.String
		purchased = append(purchased, pp)
	}

	return purchased, nil
}

// BuyProduct は商品の購入処理をトランザクションで実行します
// 💡 インターフェースや構造体を使わず、パッケージ関数として定義
func BuyProduct(productID int, buyerID string) error {
	// トランザクションの開始
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var productTitle, currentStatus string
	var sellerIDNull sql.NullString
	var price int

	// 🚀 【Claudeの神修正 1】FOR UPDATE で対象行をロックし、同時に価格(price)も取得！
	query := "SELECT title, seller_id, status, price FROM products WHERE id = ? FOR UPDATE"
	err = tx.QueryRow(query, productID).Scan(&productTitle, &sellerIDNull, &currentStatus, &price)
	if err != nil {
		return err
	}
	sellerID := sellerIDNull.String // 💡 seller_id はNULL許容カラムのため安全に受け取る

	// ステータスと不正操作のチェック
	if currentStatus == "sold" {
		return errors.New("this product is already sold")
	}
	if buyerID == sellerID {
		return errors.New("cannot buy your own product")
	}

	// 🚀 【Claudeの神修正 2】買い手のクレジット残高が足りているか厳密にチェック！
	var buyerCredits int
	err = tx.QueryRow("SELECT space_credits FROM users WHERE id = ?", buyerID).Scan(&buyerCredits)
	if err != nil {
		return err
	}
	if buyerCredits < price {
		return errors.New("insufficient credits") // 🔥 後ほどハンドラー側で捕まえます
	}

	// 🚀 トランザクション履歴の挿入（千切れていたSQLを修復）
	transactionSQL := "INSERT INTO transactions (product_id, buyer_id, seller_id, status) VALUES (?, ?, ?, 'completed')"
	if _, err = tx.Exec(transactionSQL, productID, buyerID, sellerID); err != nil {
		return err
	}

	// 🚀 【Claudeの神修正 3】お財布の移動（買い手から減算し、売り手へ加算！）
	if _, err = tx.Exec("UPDATE users SET space_credits = space_credits - ? WHERE id = ?", price, buyerID); err != nil {
		return err
	}
	if _, err = tx.Exec("UPDATE users SET space_credits = space_credits + ? WHERE id = ?", price, sellerID); err != nil {
		return err
	}

	// 商品を売り切れ状態にする
	if _, err = tx.Exec("UPDATE products SET status = 'sold' WHERE id = ?", productID); err != nil {
		return err
	}

	// 🚀 取引完了時の自動メッセージを挿入（千切れていたSQLを修復）
	autoContent := "商品【" + productTitle + "】の取引が完了しました。これよりトークルームを解放します。"
	messageSQL := "INSERT INTO messages (sender_id, receiver_id, content) VALUES (?, ?, ?)"
	if _, err = tx.Exec(messageSQL, buyerID, sellerID, autoContent); err != nil {
		return err
	}

	// すべて成功したらコミット
	return tx.Commit()
}