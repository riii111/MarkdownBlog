# Markdown Blog API

Markdown BlogのバックエンドAPIサーバーです。

## 技術スタック

- Golang
- Gin Web Framework
- Air (ホットリロード)

## プロジェクト構造

<https://github.com/golang-standards/project-layout/tree/master>を参考にしつつ、internal配下はレイヤードアーキテクチャを採用。

保守性と拡張性の高いコードベースを目指しています。

```text
backend/
├── cmd/                    # アプリケーションのエントリーポイント
│   └── api/               # APIサーバーのメイン処理
│   └── migrate/           # 手動マイグレーションのメイン処理
├── internal/              # プライベートなアプリケーションコード
│   ├── domain/           # ドメイン層
│   │   ├── model/       # ドメインモデル
│   │   └── repository/  # リポジトリ（インターフェースのみ）
│   ├── handler/         # HTTPハンドラー層
│   │   ├── dto/         # Data Transfer Objects
│   │   ├── middleware/  # HTTPミドルウェア
│   │   └── router.go    # ルーティング設定
│   │   └── user_handler.go # ユーザー関連のハンドラー
│   ├── infrastructure/ # インフラストラクチャ層
│   │   └── migration/  # マイグレーション関連の処理
│   │   └── database/   # データベース関連の処理（リポジトリの実装）
│   └── usecase/        # ユースケース層
└── pkg/                # 公開可能な再利用可能なコード
```

## アーキテクチャの説明

### レイヤー構造

1. **Domain Layer** (`internal/domain/`)
   - ビジネスロジックの中心
   - モデルとリポジトリインターフェースを定義
   - 外部依存を持たない純粋なビジネスロジック

2. **UseCase Layer** (`internal/usecase/`)
   - アプリケーションのユースケースを実装
   - ドメイン層のロジックを組み合わせてユースケースを実現
   - トランザクション管理などを担当

3. **Handler Layer** (`internal/handler/`)
   - HTTPリクエスト/レスポンスの処理
   - 入力バリデーション
   - DTOとドメインモデルの変換
   - ルーティング設定

4. **Infrastructure Layer** (`internal/infrastructure/`)
   - データベース、外部APIなどの具体的な実装
   - リポジトリインターフェースの実装
   - 外部サービスとの統合

### 依存関係の方向

```
Handler → UseCase → Domain ← Infrastructure
```

- 内側のレイヤー（Domain）は外側のレイヤーに依存しない
- 外側のレイヤーは内側のレイヤーに依存する
- Infrastructure層はDomain層のインターフェースを実装

## 開発環境のセットアップ

1. リポジトリのクローン

```bash
git clone <repository-url>
cd backend
```

2. 依存関係のインストール

```bash
go mod download
```

3. 開発サーバーの起動

```bash
air
```

サーバーは `http://localhost:8088` で起動します。

### マイグレーション

#### 1. 開発環境の場合

開発環境（`APP_ENV=dev`）では、APIサーバー起動時に自動マイグレーションが実行されます。

**制約事項：**

- 自動マイグレーションは開発環境でのみ動作します
- カラムの削除や型の変更など、破壊的な変更は自動で適用されません
- 複雑なマイグレーション（インデックスの追加など）は手動で行う必要があります

#### 2. 本番/ステージング環境の場合

本番環境やステージング環境では、手動でマイグレーションを実行する必要があります：

```bash
make migrate
```

#### マイグレーション関連コマンド

```bash
# マイグレーションの実行
make migrate

# マイグレーションの状態確認
# 各テーブルの存在有無、カラム名、データ型、NULL許可などを表示
make migrate-status

# マイグレーションのロールバック
# 注意: 開発環境でのみ実行可能。全テーブルが削除されます
make migrate-rollback
```

**注意事項：**

- ロールバックは開発環境でのみ実行可能です
- ロールバックを実行すると、全てのテーブルが削除されます
- 本番環境での使用は厳禁です
