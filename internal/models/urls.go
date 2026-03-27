package models

type URLRequest struct {
	Url         string  `json:"url"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	IsPinned    bool    `json:"is_pinned"`
	CategoryID  *int64  `json:"category_id"`
}

type ListURLsQuery struct {
	CategoryID *int64 `json:"category_id,omitempty"`
}
