package handlers

import "github.com/ramadhantriyant/buku/internal/models"

type Handler struct {
	config *models.Config
}

func New(config *models.Config) *Handler {
	return &Handler{
		config: config,
	}
}
