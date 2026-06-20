package handler

import (
	"net/http"
	"hackathon-backend/internal/model"
	"hackathon-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

// SendMessageRequest はメッセージ送信時にフロントから届くJSONの形を定義します
type SendMessageRequest struct {
	SenderID   string `json:"sender_id" binding:"required"`
	ReceiverID string `json:"receiver_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

func SendMessageHandler(c *gin.Context) {
	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := model.Message{
		SenderID:   req.SenderID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
	}

	if err := repository.SaveMessage(&msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "メッセージの保存に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, msg)
}

// GetChatHistoryHandler は2人、または特定のチャット履歴を取得します
func GetChatHistoryHandler(c *gin.Context) {
    // クエリパラメータから ID を取得 (?sender_id=xxx&receiver_id=yyy)
    senderID := c.Query("sender_id")
    receiverID := c.Query("receiver_id")

    // 必須チェック
    if senderID == "" || receiverID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "sender_id と receiver_id は必須です"})
        return
    }

    // repository 側から履歴を取得（※関数名は repository 側の実装に合わせて調整してください）
    messages, err := repository.GetChatHistory(senderID, receiverID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "履歴の取得に失敗しました"})
        return
    }

    // 履歴をそのまま JSON で返す
    c.JSON(http.StatusOK, messages)
}

func GetChatPartnersHandler(c *gin.Context) {
	// クエリパラメータかヘッダーから自分のUIDを取得
	userID := c.Query("user_id")
	if userID == "" {
		userID = c.GetHeader("X-User-UID") // どちらでも取れるようにガード
	}

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id は必須です"})
		return
	}

	partners, err := repository.GetChatPartners(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "チャット相手の取得に失敗しました"})
		return
	}

	// 相手のUIDリストをそのまま返す (例: ["mock_uid_naoya", "user_B"])
	c.JSON(http.StatusOK, partners)
}