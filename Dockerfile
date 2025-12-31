# =============================================================================
# マルチステージビルド Dockerfile
# - Stage 1: フロントエンドビルド (pnpm + Vite)
# - Stage 2: バックエンドビルド (Go)
# - Stage 3: 本番用イメージ
# =============================================================================

# -----------------------------------------------------------------------------
# Stage 1: フロントエンドビルド
# -----------------------------------------------------------------------------
FROM node:22-alpine AS frontend-builder

# pnpmをインストール
RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /app/frontend

# 依存関係のインストール（キャッシュ効率化のため先にpackage.jsonをコピー）
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# ソースコードをコピーしてビルド
COPY frontend/ ./
RUN pnpm build

# -----------------------------------------------------------------------------
# Stage 2: バックエンドビルド
# -----------------------------------------------------------------------------
FROM golang:1.25-alpine AS backend-builder

WORKDIR /app/backend

# 依存関係のインストール（キャッシュ効率化のため先にgo.mod/go.sumをコピー）
COPY backend/go.mod backend/go.sum* ./
RUN go mod download

# ソースコードをコピーしてビルド
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/server ./cmd/server

# -----------------------------------------------------------------------------
# Stage 3: 本番用イメージ
# -----------------------------------------------------------------------------
FROM alpine:3.20 AS production

# セキュリティアップデートと必要なパッケージのインストール
RUN apk --no-cache add ca-certificates tzdata

# 非rootユーザーを作成
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# バックエンドバイナリをコピー
COPY --from=backend-builder /app/server ./server

# フロントエンドのビルド成果物をコピー
COPY --from=frontend-builder /app/frontend/dist ./public

# 所有権を変更
RUN chown -R appuser:appgroup /app

# 非rootユーザーに切り替え
USER appuser

# ポートを公開
EXPOSE 8080

# ヘルスチェック
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# アプリケーションを起動
CMD ["./server"]
