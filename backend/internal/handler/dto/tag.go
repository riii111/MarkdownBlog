package dto

type TagInfo struct {
	ID   string `json:"id"` // Responseではuuid型にしない（フロントでは文字列扱いとなる）
	Name string `json:"name"`
	Slug string `json:"slug"`
}
