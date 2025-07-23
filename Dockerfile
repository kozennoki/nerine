# ビルドステージ
FROM golang:1.24-alpine AS builder

WORKDIR /app

# go mod ファイルをコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nerine cmd/server/*.go

# 最終ステージ
FROM alpine:latest

# HTTPS リクエスト用の ca-certificates をインストール
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# ビルダーステージからバイナリをコピー
COPY --from=builder /app/nerine .

# ポートを公開
EXPOSE 8080

# バイナリを実行
CMD ["./nerine"]
