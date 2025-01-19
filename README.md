# MarkdownBlog

MarkDown + Writer。
技術検証用のブログプラットフォーム。

## アプリケーション構成

[![Nuxt 3](https://img.shields.io/badge/Nuxt_3-00DC82?style=for-the-badge&logo=nuxt.js&logoColor=white)](https://nuxt.com/)
[![Vue.js](https://img.shields.io/badge/Vue.js-4FC08D?style=for-the-badge&logo=vue.js&logoColor=white)](https://vuejs.org/)
[![Nuxt UI](https://img.shields.io/badge/Nuxt_UI-00DC82?style=for-the-badge&logo=nuxt.js&logoColor=white)](https://ui.nuxt.com/)
[![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?style=for-the-badge&logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://gin-gonic.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![Traefik](https://img.shields.io/badge/Traefik-24A1C1?style=for-the-badge&logo=traefik&logoColor=white)](https://traefik.io/)

- フロントエンド: Nuxt 3 + Nuxt UI
- バックエンド: Gin(Golang)
- データベース: PostgreSQL
- プロキシ: Traefik

## ビルド〜起動まで

`make build`

`make up`　←apiコンテナのビルド完了待つ

`cd frontend`

`pnpm install`

`pnpm build`

`pnpm start`

## アクセス先

※あとで更新する↓

---

1. <http://localhost:8088/swagger/index.html>
   バックエンド（APIドキュメント）
2. <http://traefik.localhost/>
   Traefikのダッシュボード

---
