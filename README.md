# MarkdownBlog
MarkDown + Writer。
技術検証

## アプリケーション構成
[![Qwik](https://img.shields.io/badge/Qwik-4336FF?style=for-the-badge&logo=qwik&logoColor=white)](https://qwik.builder.io/)
[![Qwik City](https://img.shields.io/badge/Qwik_City-4336FF?style=for-the-badge&logo=qwik&logoColor=white)](https://qwik.builder.io/docs/qwikcity/)
[![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?style=for-the-badge&logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://gin-gonic.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![Traefik](https://img.shields.io/badge/Traefik-24A1C1?style=for-the-badge&logo=traefik&logoColor=white)](https://traefik.io/)

- フロントエンド: Qwik + Qwik City
- バックエンド: Gin(Golang)
- データベース: PostgreSQL
- プロキシ: Traefik

## ビルド〜起動まで

`make build`

`make up`　←apiコンテナのビルド完了待つ

`cd frontend`

`pnpm build`

`pnpm start`

## アクセス先

※あとで更新する↓

---

1. <http://localhost/docs/#/>
   バックエンド（APIドキュメント）
2. <http://db-admin.localhost/>
   MongoDB Express（MongoDBのGUI）
3. <http://traefik.localhost/>
   Traefikのダッシュボード

---
