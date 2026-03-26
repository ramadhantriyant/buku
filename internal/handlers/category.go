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

func (h *Handler) ListCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusForbidden, "Unauthorized")
		return
	}
	categories, err := h.config.Queries.ListCategoriesByUser(r.Context(), userID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Unable to connect to database")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, categories); err != nil {
		return
	}
}

func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	categoryID, err := strconv.ParseInt(pathParts[3], 10, 64)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, err := h.config.Queries.GetCategoryByID(r.Context(), categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSONError(w, http.StatusNotFound, "Category not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "Unable to connect to database")
		return
	}

	if category.UserID != userID {
		utils.WriteJSONError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, category); err != nil {
		return
	}
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	var req models.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Sanitize inputs
	req.Name = utils.SanitizeString(req.Name)
	req.Description = utils.SanitizeDescription(req.Description)

	if strings.TrimSpace(req.Name) == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Category name is required")
		return
	}

	params := database.CreateCategoryParams{
		Name:        req.Name,
		Description: req.Description,
		UserID:      userID,
	}

	category, err := h.config.Queries.CreateCategory(r.Context(), params)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Unable to create category")
		return
	}

	if err := utils.WriteJSON(w, http.StatusCreated, category); err != nil {
		return
	}
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	categoryID, err := strconv.ParseInt(pathParts[3], 10, 64)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	existingCategory, err := h.config.Queries.GetCategoryByID(r.Context(), categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSONError(w, http.StatusNotFound, "Category not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "Unable to connect to database")
		return
	}

	if existingCategory.UserID != userID {
		utils.WriteJSONError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	var req models.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Sanitize inputs
	req.Name = utils.SanitizeString(req.Name)
	req.Description = utils.SanitizeDescription(req.Description)

	if strings.TrimSpace(req.Name) == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Category name is required")
		return
	}

	params := database.UpdateCategoryParams{
		Name:        req.Name,
		Description: req.Description,
		ID:          categoryID,
		UserID:      userID,
	}

	category, err := h.config.Queries.UpdateCategory(r.Context(), params)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Unable to update category")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, category); err != nil {
		return
	}
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	categoryID, err := strconv.ParseInt(pathParts[3], 10, 64)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	existingCategory, err := h.config.Queries.GetCategoryByID(r.Context(), categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSONError(w, http.StatusNotFound, "Category not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "Unable to connect to database")
		return
	}

	if existingCategory.UserID != userID {
		utils.WriteJSONError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	params := database.DeleteCategoryParams{
		ID:     categoryID,
		UserID: userID,
	}

	if err := h.config.Queries.DeleteCategory(r.Context(), params); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Unable to delete category")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
