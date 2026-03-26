package models

import "time"

const (
	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	User         UserInfo `json:"user"`
}

type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RequestPasswordResetRequest struct {
	Username string `json:"username"`
}

type RequestPasswordResetResponse struct {
	Message string `json:"message"`
	// In production, this would not be included. Token is only shown here for testing.
	ResetToken string `json:"reset_token,omitempty"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}
