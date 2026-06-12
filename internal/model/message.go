package model

import "time"

// Message は messages テーブルに対応する構造体
type Message struct {
    ID                int       `json:"id"`
    SenderID          string    `json:"sender_id"`
    ReceiverID        string    `json:"receiver_id"`
    Content           string    `json:"content"`
    CreatedAt         time.Time `json:"created_at"`
}