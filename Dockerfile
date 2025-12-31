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
FROM golang:1.25 AS backend-builder

WORKDIR /app/backend

# 依存関係のインストール（キャッシュ効率化のため先にgo.mod/go.sumをコピー）
COPY backend/go.mod backend/go.sum* ./
RUN go mod download

# ソースコードをコピーしてビルド
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/server ./cmd/server

# -----------------------------------------------------------------------------
# Stage 3: 本番用イメージ (Distroless)
# -----------------------------------------------------------------------------
# Distroless: シェルやパッケージマネージャーを含まない最小イメージ
# - 攻撃対象領域を最小化
# - イメージサイズを削減
# - nonroot タグで非rootユーザーとして実行
FROM gcr.io/distroless/static-debian12:nonroot AS production

WORKDIR /app

# バックエンドバイナリをコピー（nonrootユーザー所有）
COPY --from=backend-builder --chown=nonroot:nonroot /app/server ./server

# フロントエンドのビルド成果物をコピー（nonrootユーザー所有）
COPY --from=frontend-builder --chown=nonroot:nonroot /app/frontend/dist ./public

# ポートを公開
EXPOSE 8080

# アプリケーションを起動
ENTRYPOINT ["./server"]
