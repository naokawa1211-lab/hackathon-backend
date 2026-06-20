package repository

import (
	"errors"

	"hackathon-backend/internal/db"
	"hackathon-backend/internal/model"
)

// CreateReview は取引に対するレビュー（評価）を1件登録する
// transaction_id + reviewer_id の組み合わせはUNIQUE制約があるため、二重投稿はDBエラーになる
func CreateReview(transactionID int, reviewerID, revieweeID string, rating int, comment string) error {
	if rating < 1 || rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	_, err := db.DB.Exec(
		"INSERT INTO reviews (transaction_id, reviewer_id, reviewee_id, rating, comment) VALUES (?, ?, ?, ?, ?)",
		transactionID, reviewerID, revieweeID, rating, comment,
	)
	return err
}

// GetReviewsForUser は指定したユーザーが「受け取った」レビューの一覧を新しい順で取得する
func GetReviewsForUser(userID string) ([]model.Review, error) {
	rows, err := db.DB.Query(
		"SELECT id, transaction_id, reviewer_id, reviewee_id, rating, comment, created_at FROM reviews WHERE reviewee_id = ? ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []model.Review
	for rows.Next() {
		var r model.Review
		if err := rows.Scan(&r.ID, &r.TransactionID, &r.ReviewerID, &r.RevieweeID, &r.Rating, &r.Comment, &r.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, r)
	}
	return reviews, nil
}
