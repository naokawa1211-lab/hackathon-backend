package repository

import (
	"time"
	"hackathon-backend/internal/db"
	"hackathon-backend/internal/model"
)

// 新しいメッセージをデータベースに保存します
func SaveMessage(msg *model.Message) error {
	query := `
		INSERT INTO messages (sender_id, receiver_id, content)
		VALUES (?, ?, ?)
	`
	result, err := db.DB.Exec(query, msg.SenderID, msg.ReceiverID, msg.Content)
	if err != nil {
		return err
	}

	// データベース側で自動採番されたIDを取得して構造体にセット
	id, err := result.LastInsertId()
	if err == nil {
		msg.ID = int(id)
	}

	msg.CreatedAt = time.Now()

	return nil
}

// 特定の2人のユーザー間のチャット履歴を古い順に取得
func GetChatHistory(userA, userB string) ([]model.Message, error) {
	query := `
		SELECT id, sender_id, receiver_id, content, created_at
		FROM messages
		WHERE (sender_id = ? AND receiver_id = ?)
		   OR (sender_id = ? AND receiver_id = ?)
		ORDER BY created_at ASC
	`
	rows, err := db.DB.Query(query, userA, userB, userB, userA)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []model.Message
	for rows.Next() {
		var msg model.Message
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		history = append(history, msg)
	}
	return history, nil
}