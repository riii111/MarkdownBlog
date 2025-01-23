package config

import (
	"os"
	"strconv"
)

const (
	defaultPage    = 1
	defaultPerPage = 9 // フロントエンドの表示に合わせて9件
	maxPerPage     = 100
)

// デフォルトのページ番号を取得
func GetDefaultPage() int {
	if page, err := strconv.Atoi(os.Getenv("DEFAULT_PAGE")); err == nil {
		return page
	}
	return defaultPage
}

// デフォルトの1ページあたりの件数を取得
func GetDefaultPerPage() int {
	if perPage, err := strconv.Atoi(os.Getenv("DEFAULT_PER_PAGE")); err == nil {
		return perPage
	}
	return defaultPerPage
}

// 1ページあたりの最大件数を取得
func GetMaxPerPage() int {
	if maxPerPage, err := strconv.Atoi(os.Getenv("MAX_PER_PAGE")); err == nil {
		return maxPerPage
	}
	return maxPerPage
}
