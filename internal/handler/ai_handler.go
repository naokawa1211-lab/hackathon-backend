package handler

import (
	"context"
	"fmt"
	"net/http"
	"hackathon-backend/internal/config" // 💡 ご自身のconfigのパスに合わせてください
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// フロントから送られてくるリクエストの型
type SpaceDescriptionRequest struct {
	ProductName string `json:"product_name"`
	Description string `json:"description"`
}

func GenerateSpaceDescriptionHandler(c *gin.Context) {
	var req SpaceDescriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	
	// ConfigからAPIキーを取得してGeminiクライアントを初期化
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.Env.GeminiAPIKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Geminiクライアントの初期化に失敗しました"})
		return
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-3.5-flash")

	// 💡 宇宙マーケットにふさわしいプロンプト（命令文）
	prompt := fmt.Sprintf(
		"あなたは天の川銀河で覇権を握るフリマ『MILKYWAY FLEA MARKET』の出品者です。地球人とは限らない相手への商売が天の川銀河一うまいことで名をはせています。\n"+
			"地球の商品である『%s』（元の説明: %s）を、地球人もしくは宇宙人（エイリアンやサイボーグ）に向けて、魅力的かつに紹介するユニークな商品説明文を150文字程度の日本語で作成してください。\n"+
			"SF・宇宙用語を必ず交えて、しかしいわゆる「イタく」はならないでください。中二病みたいにはならないようにお願いします。タイトルなどは不要です。本文のみを出力してください。",
		req.ProductName, req.Description,
	)

	// テキスト生成を実行
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "テキストの生成に失敗しました"})
		return
	}

	// 生成されたテキストを抽出
	var generatedText string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				generatedText += fmt.Sprintf("%v", part)
			}
		}
	}

	// フロントに返す
	c.JSON(http.StatusOK, gin.H{
		"status":            "success",
		"space_description": generatedText,
	})
}