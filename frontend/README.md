# MD Writer Frontend

Markdown記事作成・管理のためのフロントエンドアプリケーションです。

## 技術スタック

- Nuxt 3
- TypeScript
- Nuxt UI
- pnpm

## プロジェクト構造

```text
frontend/
├── components/
│   ├── common/           # 再利用可能な共通コンポーネント
│   ├── features/         # 機能単位のコンポーネント
│   └── layouts/          # レイアウトコンポーネント
├── composables/
│   ├── utils/           # 共通のcomposables
│   │   ├── useArray.ts
│   │   ├── useCookie.ts
│   │   ├── useDate.ts   # 日付フォーマットなどの共通処理
│   │   └── ...
│   ├── feature/         # 機能単位のcomposables
│   └── store/           # Pinia stores
│       ├── api/             # API関連のcomposables
│       │   └── useCoreApi.ts  # カスタムフェッチHooks
│       ├── types/               # 型定義
│       ├── constants/           # 定数
│       ├── middleware/          # ミドルウェア
│       ├── layouts/             # ページレイアウト
│       ├── pages/              # ルーティング用ページコンポーネント
│       └── stores/             # 状態管理（必要な場合）
```

## アーキテクチャの説明

### コンポーネント設計

1. **Feature-firstアプローチ**
   - 機能単位でのコンポーネント分割
   - 関心事の分離を重視
   - 再利用性と保守性の両立

2. **共通コンポーネント**
   - Nuxt UIをベースとしたUIコンポーネント
   - プロジェクト全体で再利用可能な共通部品

### API通信

`composables/api/useCoreApi.ts`を使用して、APIとの通信を一元管理：

```typescript
// 使用例
const { data } = await useCoreApi('/articles', {
  method: 'GET',
  params: { page: 1 }
})
```

## 開発環境のセットアップ

1. リポジトリのクローン

```bash
git clone <repository-url>
cd frontend
```

2. 依存関係のインストール

```bash
pnpm install
```

3. 開発サーバーの起動

```bash
pnpm dev
```

サーバーは `http://localhost:3000` で起動します。

## コーディング規約

1. **コンポーネントの命名**
   - 機能コンポーネント: `PascalCase.vue`
   - ページコンポーネント: `index.vue` または `[id].vue`

2. **Composablesの命名**
   - プレフィックスとして `use` を使用
   - キャメルケース: `useArticle.ts`

3. **型定義**
   - インターフェース: プレフィックスとして `I` を使用
   - 型エイリアス: サフィックスとして `Type` を使用

## ビルドとデプロイ

```bash
# 本番ビルド
pnpm build

# ビルドのプレビュー
pnpm preview
```

## 環境変数

`.env`ファイルで以下の環境変数を設定：

- `NUXT_PUBLIC_API_BASE`: APIのベースURL
- `NUXT_PUBLIC_API_VERSION`: APIバージョン

## テスト

```bash
# ユニットテストの実行
pnpm test

# E2Eテストの実行
pnpm test:e2e
```
