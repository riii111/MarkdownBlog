package dto

type SeriesInfo struct {
	ID    string `json:"id"` // Responseではuuid型にしない（フロントでは文字列扱いとなる）
	Title string `json:"title"`
	Slug  string `json:"slug"`
}
