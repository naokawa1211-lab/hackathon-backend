package model

import "time"

// Product は DB の products テーブルと対応する構造体です
type Product struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Price         int       `json:"price"`
	Category      string    `json:"category"`
	DurabilityPct int       `json:"durability_pct"`
	ImageURL1     string    `json:"image_url_1"`
	SellerID      string    `json:"seller_id"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type Transaction struct {
    ID        int       `json:"id"`
    ProductID int       `json:"product_id"`
    BuyerID   string    `json:"buyer_id"`
    SellerID  string    `json:"seller_id"`
    Status    string    `json:"status"`
    CreatedAt string    `json:"created_at"`
}