# Markdown Blog API

Markdown形式のブログ記事を管理するバックエンドAPI

## 技術スタック

- Golang
- Gin Web Framework
- PostgreSQL
- GORM
- Docker

## プロジェクト構成

レイヤードアーキテクチャを採用しています：

```
Handler → UseCase → Domain ← Infrastructure
```

## 開発環境のセットアップ

### 前提条件

- Docker と Docker Compose がインストールされていること
- Go の指定のVersionがインストールされていること

### 環境構築手順

1. リポジトリをクローン

```bash
git clone https://github.com/riii111/markdown-blog-api.git
cd markdown-blog-api
```

2. 環境変数ファイルの準備

```bash
cp .env.example .env
```

3. Docker Composeでサービスを起動

```bash
docker-compose up -d
```

4. マイグレーションの実行

```bash
go run cmd/migrate/main.go
```

5. APIサーバーの起動

```bash
go run cmd/api/main.go
```

または、ホットリロード機能付きで開発する場合：

```bash
air
```

## テスト

### E2Eテストの実行

E2Eテストは、実際のデータベースを使用せず、テスト用のPostgreSQLコンテナを自動的に起動して実行します。

#### ローカル開発環境での実行

簡易的なテスト実行用のシェルスクリプトを用意しています：

```bash
# 認証系APIのテストを実行
./tests/run_tests.sh
```

または、より詳細なオプションを指定してテストを実行することもできます：

```bash
# 全てのE2Eテストを実行（カバレッジ情報付き）
go test -v -cover ./tests/e2e/... -count=1

# 認証系のテストのみ実行
go test -v ./tests/e2e/auth/... -count=1
```

#### CI/CD環境での実行

GitHub Actionsでは、自動的に全てのE2Eテストが実行されます。テストの実行状況は各PRのChecksタブで確認できます。

### テストの構成

テストは以下のように構成されています（認証系の場合）

- `tests/e2e/setup.go`: テスト環境のセットアップコード
- `tests/e2e/auth/`: 認証系APIのテスト
  - `register_test.go`: ユーザー登録APIのテスト
  - `login_test.go`: ログインAPIのテスト
  - `logout_test.go`: ログアウトAPIのテスト

### テスト実装の特徴

- BDDスタイルでテストを記述
- testcontainersを使用した独立したテスト環境
- テストの分離性を確保（各テストは独立して実行可能）
- 効率的なテストヘルパーの活用
- 十分なアサーションで動作を検証

## API仕様

APIの詳細な仕様は、Swaggerドキュメントで確認できます。開発環境では以下のURLでアクセス可能です：

```
http://localhost:8080/swagger/index.html
```

## ディレクトリ構造

```
.
├── cmd/                  # エントリーポイント
│   ├── api/              # APIサーバー
│   └── migrate/          # マイグレーションツール
├── docs/                 # Swaggerドキュメント
├── internal/             # 内部パッケージ
│   ├── domain/           # ドメインモデルとリポジトリインターフェース
│   │   └── model/        # エンティティとリポジトリインターフェース
│   ├── handler/          # HTTPハンドラー
│   │   ├── dto/          # リクエスト/レスポンスのデータ構造
│   │   ├── endpoint/     # エンドポイント実装
│   │   └── middleware/   # ミドルウェア
│   ├── infrastructure/   # インフラストラクチャ層
│   │   ├── config/       # 設定
│   │   ├── database/     # データベース関連
│   │   └── migration/    # マイグレーション
│   └── usecase/          # ユースケース層
├── tests/                # テスト
│   └── e2e/              # E2Eテスト
└── traefik/              # Traefik設定
```

