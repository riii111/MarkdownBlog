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
│   ├── common/               # 再利用可能な共通コンポーネント
│   ├── features/             # 機能単位のコンポーネント
│   └── layouts/              # レイアウトコンポーネント
├── composables/
│   ├── utils/                # 共通のcomposables
│   │   ├── useArray.ts
│   │   ├── useCookie.ts
│   │   ├── useDate.ts        # 日付フォーマットなどの共通処理
│   │   └── ...
│   ├── feature/              # 機能単位のcomposables
│   ├── schemas/
│   │   ├── auth/
│   │   │   └── user/
│   │   │       ├── types.ts     # 認証関連の型定義
│   │   │       └── schema.ts    # バリデーションスキーマ
│   │   └── article/
│   │           ├── types.ts     # 記事関連の型定義
│   │           └── schema.ts    # バリデーションスキーマ
│   └── store/                   # Pinia stores
│       ├── api/                 # API関連のcomposables
│       │   └── useCoreApi.ts    # カスタムフェッチHooks
│       ├── constants/           # 定数
│       ├── middleware/          # ミドルウェア
│       ├── layouts/             # ページレイアウト
│       ├── pages/              # ルーティング用ページコンポーネント
│       └── stores/             # 状態管理（必要な場合）
```

### 1. コンポーネント設計

1. **Feature-firstアプローチ**
   - 機能単位でのコンポーネント分割
   - 関心事の分離を重視
   - 再利用性と保守性の両立

2. **共通コンポーネント**
   - Nuxt UIをベースとしたUIコンポーネント
   - プロジェクト全体で再利用可能な共通部品

#### コンポーネント実装Tips

- 検索性が上がるため、ファイル名は必ず 直前の階層構造（名前空間）に合わせたコンポーネント名にすること。

##### NG

components/
└── common/
    └── form/
         ├── Textfield.vue
         ├── Select.vue
         ├── Combobox.vue
         └── Autocomplete.vue

##### OK

components/
└── common/
    └── form/
         ├── FormTextfield.vue
         ├── FormSelect.vue
         ├── FormCombobox.vue
         └── FormAutocomplete.vue

### 2. スキーマ設計

`Single Source of Truth` を原則守る。
同じ対象のバリデーションロジックは原則一箇所で定義し、複数箇所での実装を避ける。
万が一ビジネスロジックごとに分岐する場合はusecaseを作成する。

1. **型定義とバリデーションの一元管理**
   - `composables/schemas/`に型定義とバリデーションルールを集約
   - 各機能（auth, articleなど）ごとにディレクトリを分割
   - 型の再利用性と一貫性を確保

2. **Valibot + Branded Typesの活用**
   - バリデーション済みの値を型レベルで保証
   - バンドルサイズの最適化
   - 型安全性の向上

#### スキーマ設計Tips

- バリデーション実装
  `Valibot + Branded Types`を使用して、バリデーション済みデータの型安全性を確保する

```typescript
// schemas/auth/user/types.ts
export type ValidEmail = string & { readonly brand: 'ValidEmail' };

// schemas/auth/user/schema.ts
const emailSchema = string([
  minLength(1, 'メールアドレスを入力してください'),
  email('有効なメールアドレスを入力してください')
]);

// バリデーション済みの値を作成する関数
export function createValidEmail(value: string): ValidEmail {
  const result = parse(emailSchema, value);
  return result as ValidEmail;
}

// 使用例
function sendEmail(to: ValidEmail, content: string) {
  // ここではtoが必ずバリデーション済みのメールアドレスであることが保証される
}
```

- 同じモデルのバリデーションロジックは原則一箇所で定義し、複数箇所での実装を避ける
- バリデーションエラーメッセージは、ユーザーにとってわかりやすい日本語で記述する

---

**Note:** Nuxtの自動インポート機能により、`composables/`配下のファイルは自動的にインポートされます。
明示的なインポート文は必要ありません。

### 3. API通信

`composables/api/useCoreApi.ts`を使用して、APIとの通信を一元管理：

```typescript
// 使用例
const { data } = await useCoreApi('/articles', {
  method: 'GET',
  params: { page: 1 }
})
```

#### API型定義のTips

- リクエスト/レスポンスの型定義には、それぞれ`Request`/`Response`というサフィックスを使用

  ```typescript
  interface IUserRequest { ... }  // リクエストボディの型
  interface IUserResponse { ... } // レスポンスボディの型
  ```

- この命名規則により、データの流れが型レベルで明確になり、コードの可読性が向上する

### 4. 定数管理

定数の管理は以下の基準に従って判断します：

1. **グローバル定数（constants/配下に配置）**
   - アプリケーション全体で共有される設定値

   ```typescript
   // constants/auth.ts
   export const AUTH_CONSTANTS = {
     TOKEN_KEY: "auth_token",
     API_ENDPOINT: "/api/v1/auth",
   } as const;
   ```

2. **機能スコープの定数（各機能のファイル内で定義）**
   - 特定の機能に閉じた値やメッセージ

   ```typescript
   // schemas/auth/user/schema.ts
   const VALIDATION_RULES = {
     PASSWORD: {
       MIN_LENGTH: 8,
       MAX_LENGTH: 100,
     },
     // ... その他のルール
   } as const;
   ```

#### 定数管理のTips

- バリデーションメッセージやルールは、使用される場所の近くで定義する
- 国際化（i18n）が必要な文言は必ず定数化する
- 数値の閾値は、変更の可能性を考慮して定数化を検討する
- 定数のスコープは可能な限り小さく保つ

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
