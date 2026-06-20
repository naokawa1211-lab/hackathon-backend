package storage

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"hackathon-backend/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var r2Client *s3.Client

// InitR2 はCloudflare R2(S3互換)向けのクライアントを初期化します
func InitR2() error {
	cfg := config.Env

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.R2AccountID)

	awsCfg := aws.Config{
		Region:      "auto",
		Credentials: credentials.NewStaticCredentialsProvider(cfg.R2AccessKeyID, cfg.R2SecretAccessKey, ""),
	}

	r2Client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return nil
}

// UploadImage は受け取った画像ファイルをR2バケットにアップロードし、公開URLを返します
func UploadImage(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	key := fmt.Sprintf("products/%d%s", time.Now().UnixNano(), ext)

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = r2Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(config.Env.R2BucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Printf("[R2] PutObject failed: bucket=%q key=%q err=%v", config.Env.R2BucketName, key, err)
		return "", err
	}

	publicURL := strings.TrimRight(config.Env.R2PublicURL, "/") + "/" + key
	return publicURL, nil
}
