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
├── internal/              # プライベートなアプリケーションコード
│   ├── domain/           # ドメイン層
│   │   ├── model/       # ドメインモデル
│   │   └── repository/  # リポジトリインターフェース
│   ├── handler/         # HTTPハンドラー層
│   │   ├── dto/        # Data Transfer Objects
│   │   ├── middleware/ # HTTPミドルウェア
│   │   └── router/     # ルーティング設定
│   ├── infrastructure/ # インフラストラクチャ層
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
