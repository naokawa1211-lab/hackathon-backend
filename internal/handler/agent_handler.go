package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"hackathon-backend/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// AgentProductContext はPolarisに渡す商品情報（鑑定対象 / コンシェルジュ候補リスト共通）
type AgentProductContext struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Category    string `json:"category"`
}

// AgentRequest はフロントから届く統合AIエージェント「Polaris」へのリクエスト
type AgentRequest struct {
	Mode     string                 `json:"mode" binding:"required"`    // "appraise" or "concierge"
	Message  string                 `json:"message" binding:"required"` // ユーザーの入力文
	Product  *AgentProductContext   `json:"product,omitempty"`          // 鑑定モード: 鑑定対象の商品
	Products []AgentProductContext  `json:"products,omitempty"`         // コンシェルジュモード: 提案候補の商品リスト
}

const (
	agentModeAppraise  = "appraise"
	agentModeConcierge = "concierge"
)

// AgentHandler は統合AIエージェント「Polaris」の応答を生成します
// mode に応じてSystem Instructionを切り替え、Geminiに1回投げるだけのシンプルな実装
func AgentHandler(c *gin.Context) {
	var req AgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mode と message は必須です"})
		return
	}

	var systemInstruction, contextBlock string

	switch req.Mode {
	case agentModeAppraise:
		systemInstruction = "あなた（Polaris）は銀河一のガジェット鑑定士です。送られた商品データ（地球のiPadのような未知のテクノロジーや宇宙のジャンク品）を、" +
			"それを全く知らない異星人に向けて、ユーモアを交えつつ「何に使うものか」「良品・不良品の見分け方」「宇宙空間での使用上の注意」を噛み砕いて説明してください。"
		if req.Product != nil {
			contextBlock = fmt.Sprintf(
				"【鑑定対象の商品データ】\nタイトル: %s\n説明文: %s\n価格: %d円\nカテゴリ: %s\n",
				req.Product.Title, req.Product.Description, req.Product.Price, req.Product.Category,
			)
		}
	case agentModeConcierge:
		systemInstruction = "あなた（Polaris）は宇宙フリマの案内人です。ユーザーの要望に合う商品を、" +
			"今あるリスト（フロントから渡されるコンテキスト）の中からユーモアたっぷりに提案してください。"
		if len(req.Products) > 0 {
			var sb strings.Builder
			sb.WriteString("【現在出品中の商品リスト】\n")
			for _, p := range req.Products {
				sb.WriteString(fmt.Sprintf("- 「%s」 / %d円 / カテゴリ: %s\n", p.Title, p.Price, p.Category))
			}
			contextBlock = sb.String()
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "mode は appraise または concierge を指定してください"})
		return
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.Env.GeminiAPIKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Geminiクライアントの初期化に失敗しました"})
		return
	}
	defer client.Close()

	model := client.GenerativeModel(geminiModelName)
	model.SystemInstruction = &genai.Content{Parts: []genai.Part{genai.Text(systemInstruction)}}

	prompt := req.Message
	if contextBlock != "" {
		prompt = contextBlock + "\n【ユーザーからの発言】\n" + req.Message
	}

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Polarisとの通信に失敗しました"})
		return
	}

	var reply string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				reply += fmt.Sprintf("%v", part)
			}
		}
	}

	if strings.TrimSpace(reply) == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Polarisから有効な応答が得られませんでした"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reply": reply})
}
