.PHONY: genswag genproto run gomigrate migrateup migratedown migrateforce migrateversion goupdate

# DB接続設定（環境変数で上書き可能）
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= todoapp
DB_SSLMODE ?= disable
DATABASE_URL ?= postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

run:
	go run cmd/app/main.go

genswag:
	protoc -I . --openapiv2_out ./docs --openapiv2_opt allow_merge=true,disable_default_errors=true $(file)

genproto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/v1/*.proto

# マイグレーションファイル作成: make gomigrate file=create_users
gomigrate:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate create -ext sql -dir db/migrations -seq $(file)

# マイグレーション実行（全て適用）
migrateup:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate -path db/migrations -database '$(DATABASE_URL)' -verbose up

# マイグレーション実行（指定数だけ適用）: make migrateup-n n=1
migrateup-n:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate -path db/migrations -database '$(DATABASE_URL)' -verbose up $(n)

# マイグレーションロールバック（全て）
migratedown:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate -path db/migrations -database '$(DATABASE_URL)' -verbose down

# マイグレーションロールバック（指定数）: make migratedown-n n=1
migratedown-n:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate -path db/migrations -database '$(DATABASE_URL)' -verbose down $(n)

# マイグレーション強制バージョン設定: make migrateforce v=3
migrateforce:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate -path db/migrations -database '$(DATABASE_URL)' force $(v)

# マイグレーションバージョン確認
migrateversion:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate -path db/migrations -database '$(DATABASE_URL)' version

goupdate:
	go get -t -u ./...
	go mod tidy