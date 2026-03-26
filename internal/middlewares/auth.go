package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ramadhantriyant/buku/internal/utils"
)

type contextKey string

const (
	contextKeyUserID  contextKey = "userID"
	contextKeyIsAdmin contextKey = "isAdmin"
)

type Claims struct {
	UserID  int64  `json:"user_id"`
	Type    string `json:"type"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(contextKeyUserID).(int64)
	return userID, ok
}

func GetIsAdminFromContext(ctx context.Context) bool {
	isAdmin, ok := ctx.Value(contextKeyIsAdmin).(bool)
	return ok && isAdmin
}

func GenerateToken(userID int64, jwtSecret string, tokenType string, duration time.Duration, isAdmin bool) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:  userID,
		Type:    tokenType,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ParseToken(tokenString string, jwtSecret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func AuthMiddleware(jwtSecret string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.WriteJSONError(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				utils.WriteJSONError(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			tokenString := parts[1]
			claims, err := ParseToken(tokenString, jwtSecret)
			if err != nil {
				utils.WriteJSONError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			if claims.Type != "access" {
				utils.WriteJSONError(w, http.StatusUnauthorized, "Invalid token type")
				return
			}

			ctx := context.WithValue(r.Context(), contextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, contextKeyIsAdmin, claims.IsAdmin)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAdmin := GetIsAdminFromContext(r.Context())
		if !isAdmin {
			utils.WriteJSONError(w, http.StatusForbidden, "Admin access required")
			return
		}
		next.ServeHTTP(w, r)
	}
}
