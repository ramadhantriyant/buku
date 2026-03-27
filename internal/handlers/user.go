package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/ramadhantriyant/buku/internal/database"
	"github.com/ramadhantriyant/buku/internal/middlewares"
	"github.com/ramadhantriyant/buku/internal/models"
	"github.com/ramadhantriyant/buku/internal/utils"
)

const (
	accessTokenExpiry  = 15 * time.Minute
	refreshTokenExpiry = 7 * 24 * time.Hour
)

func hashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func checkPasswordHash(password, hash string) bool {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	return err == nil && match
}

func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Sanitize inputs
	req.Username = utils.SanitizeString(req.Username)
	req.Name = utils.SanitizeString(req.Name)

	if strings.TrimSpace(req.Username) == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Username is required")
		return
	}

	if strings.TrimSpace(req.Password) == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Password is required")
		return
	}

	if len(req.Password) < 6 {
		utils.WriteJSONError(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Name is required")
		return
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	// Check if this is the first user (make them admin)
	userCount, err := h.config.Queries.CountUsers(r.Context())
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to check user count")
		return
	}
	isFirstUser := userCount == 0

	params := database.CreateUserParams{
		Username: req.Username,
		Password: hashedPassword,
		Name:     req.Name,
		IsAdmin:  false,
	}
	if isFirstUser {
		params.IsAdmin = true
	}

	user, err := h.config.Queries.CreateUser(r.Context(), params)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			utils.WriteJSONError(w, http.StatusConflict, "Username already exists")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	accessToken, err := middlewares.GenerateToken(user.ID, h.config.JWTSecret, "access", accessTokenExpiry, user.IsAdmin)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	refreshTokenString, err := generateRefreshToken()
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	refreshTokenParams := database.CreateRefreshTokenParams{
		UserID:    user.ID,
		TokenHash: refreshTokenString,
		ExpiresAt: time.Now().Add(refreshTokenExpiry),
	}

	if _, err := h.config.Queries.CreateRefreshToken(r.Context(), refreshTokenParams); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to save refresh token")
		return
	}

	response := models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(accessTokenExpiry.Seconds()),
		User: models.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
		},
	}

	if err := utils.WriteJSON(w, http.StatusCreated, response); err != nil {
		return
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	user, err := h.config.Queries.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSONError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to authenticate")
		return
	}

	if !checkPasswordHash(req.Password, user.Password) {
		utils.WriteJSONError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	accessToken, err := middlewares.GenerateToken(user.ID, h.config.JWTSecret, "access", accessTokenExpiry, user.IsAdmin)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	refreshTokenString, err := generateRefreshToken()
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	refreshTokenParams := database.CreateRefreshTokenParams{
		UserID:    user.ID,
		TokenHash: refreshTokenString,
		ExpiresAt: time.Now().Add(refreshTokenExpiry),
	}

	if _, err := h.config.Queries.CreateRefreshToken(r.Context(), refreshTokenParams); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to save refresh token")
		return
	}

	response := models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(accessTokenExpiry.Seconds()),
		User: models.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
		},
	}

	if err := utils.WriteJSON(w, http.StatusOK, response); err != nil {
		return
	}
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	refreshToken, err := h.config.Queries.GetActiveRefreshTokenByHash(r.Context(), req.RefreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSONError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to validate refresh token")
		return
	}

	user, err := h.config.Queries.GetUserByID(r.Context(), refreshToken.UserID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	accessToken, err := middlewares.GenerateToken(user.ID, h.config.JWTSecret, "access", accessTokenExpiry, user.IsAdmin)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	newRefreshTokenString, err := generateRefreshToken()
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	if err := h.config.Queries.RevokeRefreshToken(r.Context(), req.RefreshToken); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to revoke old refresh token")
		return
	}

	refreshTokenParams := database.CreateRefreshTokenParams{
		UserID:    user.ID,
		TokenHash: newRefreshTokenString,
		ExpiresAt: time.Now().Add(refreshTokenExpiry),
	}

	if _, err := h.config.Queries.CreateRefreshToken(r.Context(), refreshTokenParams); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to save refresh token")
		return
	}

	response := models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshTokenString,
		ExpiresIn:    int64(accessTokenExpiry.Seconds()),
		User: models.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
		},
	}

	if err := utils.WriteJSON(w, http.StatusOK, response); err != nil {
		return
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req models.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	if err := h.config.Queries.RevokeRefreshToken(r.Context(), req.RefreshToken); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.config.Queries.GetUserByID(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSONError(w, http.StatusNotFound, "User not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	response := models.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
	}

	if err := utils.WriteJSON(w, http.StatusOK, response); err != nil {
		return
	}
}

func generatePasswordResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (h *Handler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req models.RequestPasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.Username) == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Username is required")
		return
	}

	// Check if user exists
	user, err := h.config.Queries.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return success even if user doesn't exist (security best practice)
			// This prevents user enumeration attacks
			response := models.RequestPasswordResetResponse{
				Message: "If the user exists, a password reset link has been sent",
			}
			if err := utils.WriteJSON(w, http.StatusOK, response); err != nil {
				return
			}
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to process request")
		return
	}

	// Delete any existing password reset tokens for this user
	if err := h.config.Queries.DeleteUserPasswordResetTokens(r.Context(), user.ID); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to process request")
		return
	}

	// Generate new token
	token, err := generatePasswordResetToken()
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to generate reset token")
		return
	}

	// Save token to database (expires in 1 hour)
	params := database.CreatePasswordResetTokenParams{
		UserID:    user.ID,
		TokenHash: token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if _, err := h.config.Queries.CreatePasswordResetToken(r.Context(), params); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to save reset token")
		return
	}

	// In production, you would send an email here with the reset link
	// For development/testing, we return the token in the response
	response := models.RequestPasswordResetResponse{
		Message:    "If the user exists, a password reset link has been sent",
		ResetToken: token, // Remove this in production!
	}

	if err := utils.WriteJSON(w, http.StatusOK, response); err != nil {
		return
	}
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.Token) == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Token is required")
		return
	}

	if strings.TrimSpace(req.Password) == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Password is required")
		return
	}

	if len(req.Password) < 6 {
		utils.WriteJSONError(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	// Get the password reset token
	resetToken, err := h.config.Queries.GetPasswordResetTokenByHash(r.Context(), req.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSONError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to validate token")
		return
	}

	// Hash the new password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	// Update user's password
	updateParams := database.UpdateUserPasswordParams{
		Password: hashedPassword,
		ID:       resetToken.UserID,
	}

	if err := h.config.Queries.UpdateUserPassword(r.Context(), updateParams); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to update password")
		return
	}

	// Mark token as used
	if err := h.config.Queries.MarkPasswordResetTokenUsed(r.Context(), resetToken.ID); err != nil {
		// Log error but don't fail the request - password was already updated
		// In production, you would log this error
	}

	// Revoke all refresh tokens for this user (force re-login)
	if err := h.config.Queries.RevokeAllUserRefreshTokens(r.Context(), resetToken.UserID); err != nil {
		// Log error but don't fail the request
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.CurrentPassword) == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Current password is required")
		return
	}

	if strings.TrimSpace(req.NewPassword) == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "New password is required")
		return
	}

	if len(req.NewPassword) < 6 {
		utils.WriteJSONError(w, http.StatusBadRequest, "New password must be at least 6 characters")
		return
	}

	user, err := h.config.Queries.GetUserByID(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSONError(w, http.StatusNotFound, "User not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	if !checkPasswordHash(req.CurrentPassword, user.Password) {
		utils.WriteJSONError(w, http.StatusUnauthorized, "Current password is incorrect")
		return
	}

	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	if err := h.config.Queries.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
		Password: hashedPassword,
		ID:       userID,
	}); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to update password")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Admin dashboard - admin access granted",
		"admin":   true,
	}

	if err := utils.WriteJSON(w, http.StatusOK, response); err != nil {
		return
	}
}
