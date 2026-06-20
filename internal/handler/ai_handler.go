package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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

	model := client.GenerativeModel(geminiModelName)

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

// Geminiのマルチモーダル(画像入力)に対応したモデル名
const geminiModelName = "gemini-2.5-flash"

// AnalyzeImageResponse はGeminiの画像認識結果からフロントへ返すJSONの形
type AnalyzeImageResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// AnalyzeImageHandler は出品画像をGeminiのマルチモーダルAPIに渡し、
// 画像に写っているものをベースにしたSF風の商品タイトル・説明文を生成します
func AnalyzeImageHandler(c *gin.Context) {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image ファイルは必須です"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "画像の読み込みに失敗しました"})
		return
	}
	defer file.Close()

	imgBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "画像の読み込みに失敗しました"})
		return
	}

	imageFormat := imageFormatFromContentType(fileHeader.Header.Get("Content-Type"))

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.Env.GeminiAPIKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Geminiクライアントの初期化に失敗しました"})
		return
	}
	defer client.Close()

	model := client.GenerativeModel(geminiModelName)

	prompt := `添付された画像に写っているものを識別してください。その上で、このアプリ「Milkyway Flea Market（宇宙のフリマアプリ）」に出品するための、SF・宇宙パロディ風の魅力的な「商品タイトル」と「商品説明文」を日本語で生成してください。
レスポンスは前置きやコードブロック記号(` + "```" + `)を一切付けず、以下のJSON形式のみで返却してください。
{"title": "宇宙風のタイトル", "description": "宇宙風の説明文"}`

	resp, err := model.GenerateContent(ctx, genai.ImageData(imageFormat, imgBytes), genai.Text(prompt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "画像解析に失敗しました"})
		return
	}

	var rawText string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				rawText += fmt.Sprintf("%v", part)
			}
		}
	}

	if strings.TrimSpace(rawText) == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Geminiから有効な応答が得られませんでした"})
		return
	}

	var result AnalyzeImageResponse
	if err := json.Unmarshal([]byte(stripJSONFence(rawText)), &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AIレスポンスの解析に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// stripJSONFence はGeminiがレスポンスをMarkdownのコードブロックで囲んで返してくる場合に備えて除去する
func stripJSONFence(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}

// imageFormatFromContentType はContent-Typeから genai.ImageData が要求する形式文字列("jpeg"/"png"等)を取り出す
func imageFormatFromContentType(contentType string) string {
	switch {
	case strings.Contains(contentType, "png"):
		return "png"
	case strings.Contains(contentType, "webp"):
		return "webp"
	case strings.Contains(contentType, "heic"):
		return "heic"
	case strings.Contains(contentType, "gif"):
		return "gif"
	default:
		return "jpeg"
	}
}
