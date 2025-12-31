# TODO API Backend

Clean ArchitectureパターンでGolangで実装されたTODOアプリケーションのバックエンドです。

## アーキテクチャ

このプロジェクトはClean Architectureの原則に従って設計されています：

```
backend/
├── cmd/
│   └── server/          # エントリーポイント
│       └── main.go
├── internal/
│   ├── model/           # 共有データ構造
│   │   ├── todo.go
│   │   └── error.go
│   ├── domain/          # ビジネスロジック
│   │   └── repository/  # リポジトリインターフェース
│   │       └── todo.go
│   ├── application/     # ユースケース
│   │   └── usecase/
│   │       └── todo.go
│   ├── infrastructure/  # 外部依存の実装
│   │   └── persistence/
│   │       ├── database.go
│   │       └── todo_repository.go
│   ├── interface/       # HTTPハンドラー
│   │   └── handler/
│   │       └── todo_handler.go
│   └── router/          # ルーティング
│       └── router.go
├── go.mod
├── go.sum
└── .env.example
```

### 各層の責務

- **model**: ドメインとアプリケーション層で共有されるデータ構造とエラー定義
- **domain**: ビジネスロジックとリポジトリインターフェース
- **application**: ユースケースの実装（ビジネスロジックの調整）
- **infrastructure**: 外部依存の実装（データベース、API等）
- **interface**: HTTPハンドラーの実装
- **router**: ルーティング設定とミドルウェア
- **cmd**: アプリケーションのエントリーポイント

## 機能

- TODO作成（Create）
- TODO取得（Read）
  - 単一TODO取得
  - 全TODO取得
- TODO更新（Update）
- TODO削除（Delete）

## API仕様

### エンドポイント

| メソッド | パス | 説明 |
|---------|------|------|
| POST | /api/v1/todos | TODOを作成 |
| GET | /api/v1/todos | 全TODOを取得 |
| GET | /api/v1/todos/{id} | 指定IDのTODOを取得 |
| PUT | /api/v1/todos/{id} | 指定IDのTODOを更新 |
| DELETE | /api/v1/todos/{id} | 指定IDのTODOを削除 |
| GET | /health | ヘルスチェック |

### リクエスト例

#### TODO作成
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "買い物に行く",
    "description": "牛乳とパンを買う"
  }'
```

#### TODO更新
```bash
curl -X PUT http://localhost:8080/api/v1/todos/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "title": "買い物完了",
    "completed": true
  }'
```

### レスポンス形式

成功時はTODOオブジェクトを返します：
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "買い物に行く",
  "description": "牛乳とパンを買う",
  "completed": false,
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

エラー時はRFC 9457に準拠したProblem Details形式を返します：
```json
{
  "type": "about:blank",
  "title": "Not Found",
  "status": 404,
  "detail": "指定されたTODOが見つかりません",
  "instance": "/api/v1/todos/123"
}
```

## セットアップ

### 前提条件

- Go 1.25以上
- PostgreSQL 12以上

### データベースのセットアップ

1. PostgreSQLデータベースを作成：
```bash
createdb todoapp
```

2. 環境変数を設定：
```bash
cp .env.example .env
# .envファイルを編集してデータベース接続情報を設定
```

### 依存関係のインストール

```bash
cd backend
go mod download
```

### 起動

```bash
go run cmd/server/main.go
```

サーバーは`http://localhost:8080`で起動します。

### ビルド

```bash
go build -o bin/server cmd/server/main.go
./bin/server
```

## 技術スタック

- **言語**: Go 1.25+
- **Webフレームワーク**: gorilla/mux
- **データベース**: PostgreSQL
- **ドライバー**: lib/pq
- **UUID生成**: google/uuid
- **CORS**: rs/cors
- **ロギング**: log/slog (標準ライブラリ)

## 主要な設計パターン

### Clean Architecture

依存関係は外側から内側に向かい、ビジネスロジックは外部の詳細から独立しています：
- domain層は他の層に依存しない
- application層はdomainに依存
- infrastructure層はdomainとapplicationに依存
- interface層はapplicationに依存

### Dependency Injection

依存性の注入により、テストが容易で疎結合な設計になっています。

### Repository Pattern

データアクセスを抽象化し、ビジネスロジックをデータソースから分離しています。

### Error Handling

- `panic`は使用せず、すべてのエラーは`error`型で返される
- エラーは`fmt.Errorf`でラップしてコンテキストを追加
- 共通エラーは`internal/model/error.go`で定義

### Structured Logging

- `log/slog`を使用した構造化ログ
- コンテキストを通じてログ情報を伝播
- 適切なログレベル（Info, Warn, Error）の使い分け

### Graceful Shutdown

シグナル受信時にグレースフルシャットダウンを実行し、処理中のリクエストを完了させます。

## 環境変数

| 変数名 | 説明 | デフォルト値 |
|-------|------|------------|
| DB_HOST | データベースホスト | localhost |
| DB_PORT | データベースポート | 5432 |
| DB_USER | データベースユーザー | postgres |
| DB_PASSWORD | データベースパスワード | postgres |
| DB_NAME | データベース名 | todoapp |
| DB_SSLMODE | SSL接続モード | disable |
| PORT | サーバーポート | 8080 |

## 開発

### コーディング規約

プロジェクトルールに従ってください：
- 可読性優先
- 構造化ログの使用（log/slog）
- コンテキスト伝播
- Panic禁止
- エラーラップ

詳細は`agents.md`を参照してください。

## ライセンス

MIT License
