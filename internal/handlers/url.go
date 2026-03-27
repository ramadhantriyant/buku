package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ramadhantriyant/buku/internal/database"
	"github.com/ramadhantriyant/buku/internal/middlewares"
	"github.com/ramadhantriyant/buku/internal/models"
	"github.com/ramadhantriyant/buku/internal/utils"
)

func (h *Handler) ListURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	var urls []database.Url
	var err error

	// Check for search query
	searchQuery := r.URL.Query().Get("search")
	if searchQuery != "" {
		// Sanitize search query
		searchQuery = utils.SanitizeString(searchQuery)
		// Add wildcards for LIKE query
		searchPattern := "%" + searchQuery + "%"
		params := database.SearchURLsParams{
			UserID:      userID,
			Url:         searchPattern,
			Title:       &searchPattern,
			Description: &searchPattern,
		}
		urls, err = h.config.Queries.SearchURLs(r.Context(), params)
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Unable to search URLs")
			return
		}
	} else if categoryIDStr := r.URL.Query().Get("category_id"); categoryIDStr != "" {
		categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
		if err != nil {
			utils.WriteJSONError(w, http.StatusBadRequest, "Invalid category ID")
			return
		}
		params := database.ListURLsByCategoryParams{
			CategoryID: &categoryID,
			UserID:     userID,
		}
		urls, err = h.config.Queries.ListURLsByCategory(r.Context(), params)
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Unable to fetch URLs")
			return
		}
	} else {
		urls, err = h.config.Queries.ListURLsByUser(r.Context(), userID)
		if err != nil {
			utils.WriteJSONError(w, http.StatusInternalServerError, "Unable to fetch URLs")
			return
		}
	}

	if err := utils.WriteJSON(w, http.StatusOK, urls); err != nil {
		return
	}
}

func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid URL ID")
		return
	}

	urlID, err := strconv.ParseInt(pathParts[3], 10, 64)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid URL ID")
		return
	}

	url, err := h.config.Queries.GetURLByID(r.Context(), urlID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSONError(w, http.StatusNotFound, "URL not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "Unable to fetch URL")
		return
	}

	if url.UserID != userID {
		utils.WriteJSONError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, url); err != nil {
		return
	}
}

func (h *Handler) CreateURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	var req models.URLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Sanitize inputs
	req.Url = utils.SanitizeURL(req.Url)
	req.Description = utils.SanitizeDescription(req.Description)

	if strings.TrimSpace(req.Url) == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "URL is required")
		return
	}

	// Validate category ownership if provided
	if req.CategoryID != nil {
		cat, err := h.config.Queries.GetCategoryByID(r.Context(), *req.CategoryID)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.WriteJSONError(w, http.StatusBadRequest, "Category not found")
				return
			}
			utils.WriteJSONError(w, http.StatusInternalServerError, "Unable to validate category")
			return
		}
		if cat.UserID != userID {
			utils.WriteJSONError(w, http.StatusForbidden, "Category does not belong to you")
			return
		}
	}

	params := database.CreateURLParams{
		Url:         req.Url,
		Title:       req.Title,
		Description: req.Description,
		IsPinned:    req.IsPinned,
		CategoryID:  req.CategoryID,
		UserID:      userID,
	}

	url, err := h.config.Queries.CreateURL(r.Context(), params)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Unable to create URL")
		return
	}

	if err := utils.WriteJSON(w, http.StatusCreated, url); err != nil {
		return
	}
}

func (h *Handler) UpdateURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid URL ID")
		return
	}

	urlID, err := strconv.ParseInt(pathParts[3], 10, 64)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid URL ID")
		return
	}

	existingURL, err := h.config.Queries.GetURLByID(r.Context(), urlID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSONError(w, http.StatusNotFound, "URL not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "Unable to fetch URL")
		return
	}

	if existingURL.UserID != userID {
		utils.WriteJSONError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	var req models.URLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Sanitize inputs
	req.Url = utils.SanitizeURL(req.Url)
	req.Description = utils.SanitizeDescription(req.Description)

	if strings.TrimSpace(req.Url) == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "URL is required")
		return
	}

	// Validate category ownership if provided
	if req.CategoryID != nil {
		cat, err := h.config.Queries.GetCategoryByID(r.Context(), *req.CategoryID)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.WriteJSONError(w, http.StatusBadRequest, "Category not found")
				return
			}
			utils.WriteJSONError(w, http.StatusInternalServerError, "Unable to validate category")
			return
		}
		if cat.UserID != userID {
			utils.WriteJSONError(w, http.StatusForbidden, "Category does not belong to you")
			return
		}
	}

	params := database.UpdateURLParams{
		Url:         req.Url,
		Title:       req.Title,
		Description: req.Description,
		IsPinned:    req.IsPinned,
		CategoryID:  req.CategoryID,
		ID:          urlID,
		UserID:      userID,
	}

	url, err := h.config.Queries.UpdateURL(r.Context(), params)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Unable to update URL")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, url); err != nil {
		return
	}
}

func (h *Handler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid URL ID")
		return
	}

	urlID, err := strconv.ParseInt(pathParts[3], 10, 64)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid URL ID")
		return
	}

	existingURL, err := h.config.Queries.GetURLByID(r.Context(), urlID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSONError(w, http.StatusNotFound, "URL not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "Unable to fetch URL")
		return
	}

	if existingURL.UserID != userID {
		utils.WriteJSONError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	params := database.DeleteURLParams{
		ID:     urlID,
		UserID: userID,
	}

	if err := h.config.Queries.DeleteURL(r.Context(), params); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Unable to delete URL")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
