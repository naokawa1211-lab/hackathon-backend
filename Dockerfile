# ステージ1: Goのコードをビルドして実行ファイルを作る環境
FROM golang:1.22-alpine AS builder
WORKDIR /app

# 依存関係（go.mod / go.sum）をコピーしてダウンロード
COPY go.mod go.sum ./
RUN go mod download

# すべてのソースコードをコピーしてコンパイル
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

# ステージ2: 実行ファイルだけを動かす
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app

# ステージ1で作った実行ファイルだけを持ってくる
COPY --from=builder /app/main .

# Cloud Runが使うポートを開放
EXPOSE 8080

# アプリを起動！
CMD ["./main"]