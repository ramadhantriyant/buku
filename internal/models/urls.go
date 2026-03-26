package models

type URLRequest struct {
	Url         string  `json:"url"`
	Description *string `json:"description"`
	CategoryID  int64   `json:"category_id"`
}

type ListURLsQuery struct {
	CategoryID *int64 `json:"category_id,omitempty"`
}
