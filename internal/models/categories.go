package models

type CategoryRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}
