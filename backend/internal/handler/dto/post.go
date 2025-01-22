package dto

type CreatePostRequest struct {
	Title    string   `json:"title" binding:"required,max=255"`
	Content  string   `json:"content" binding:"required"`
	SeriesID *string  `json:"series_id,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

type CreatePostResponse struct {
	Slug string `json:"slug"`
}
