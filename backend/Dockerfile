# =============================================================================
# バックエンド用 Dockerfile (Railway デプロイ用)
# =============================================================================

# -----------------------------------------------------------------------------
# Stage 1: ビルドステージ
# -----------------------------------------------------------------------------
FROM golang:1.23-alpine AS builder

# 必要なパッケージをインストール
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# 依存関係のインストール（キャッシュ効率化のため先にgo.mod/go.sumをコピー）
COPY go.mod go.sum* ./
RUN go mod download

# ソースコードをコピーしてビルド
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/server ./cmd/server

# -----------------------------------------------------------------------------
# Stage 2: 本番用イメージ
# -----------------------------------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

# CA証明書をコピー（HTTPS通信用）
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# バイナリをコピー
COPY --from=builder --chown=nonroot:nonroot /app/server ./server

# Railwayは環境変数でPORTを指定する
# EXPOSE $PORT

# アプリケーションを起動
ENTRYPOINT ["./server"]
