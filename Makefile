down:
	docker compose down -v

network:
	docker network create devtrackr-network

prune:
	docker system prune --force

build:
	docker compose build

build-n:
	docker compose build --no-cache

up:
	docker compose up -d

re:
	make down
	make build
	make up

re-n:
	make down
	make build-n
	make up

re-all-force:
	make down
	make prune
	make network
	make build-n
	make up

container=frontend
command=/bin/bash
tail=200

logs:
	docker compose logs --tail=${tail} ${container}

logs-all:
	docker compose logs --tail=${tail}

exec:
	docker compose exec ${container} ${command}

# リンターとフォーマッター関連のコマンド
lint:
	docker compose exec api golangci-lint run ./...

fmt:
	docker compose exec api go fmt ./...

# コードの品質チェック（フォーマット＋リント）
check:
	make fmt
	make lint

# テスト実行
test:
	docker compose exec api go test ./...

# 全ての検証を実行（フォーマット、リント、テスト）
verify:
	make fmt
	make lint
	make test

# コードジェネレート
generate:
	docker compose exec api go generate ./...

# マイグレーション関連のコマンド
.PHONY: migrate migrate-status migrate-rollback

# マイグレーション実行
migrate:
	docker compose exec api go run cmd/migrate/main.go \
		-host=$${POSTGRES_HOST:-db} \
		-port=$${POSTGRES_PORT:-5432} \
		-dbname=$${POSTGRES_DB:-markdown-blog} \
		-user=$${POSTGRES_USER:-root} \
		-password=$${POSTGRES_PASSWORD:-password}

# マイグレーションのステータス確認（必要に応じて実装）
migrate-status:
	docker compose exec api go run cmd/migrate/main.go -status

# マイグレーションのロールバック（必要に応じて実装）
migrate-rollback:
	docker compose exec api go run cmd/migrate/main.go -rollback