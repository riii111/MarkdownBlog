package main

import (
	"log"
	"net/http"

	"github.com/riii111/markdown-blog-api/internal/health"
)

func main() {
	// ヘルスチェックエンドポイントの登録
	http.HandleFunc("/health", health.Handler)

	// サーバーの起動
	log.Println("Server starting on :8088")
	if err := http.ListenAndServe(":8088", nil); err != nil {
		log.Fatal(err)
	}
}
