package handlers

import (
    "context"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDKey contextKey = "user_id"

// AuthMiddleware kiểm tra JWT và thêm user_id vào context
func (ah *AuthHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing token", http.StatusUnauthorized)
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return ah.jwtSecret, nil
        })
        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            http.Error(w, "Invalid token claims", http.StatusUnauthorized)
            return
        }

        userID, ok := claims["user_id"].(string)
        if !ok {
            http.Error(w, "Invalid user ID", http.StatusUnauthorized)
            return
        }

        // Thêm user_id vào context
        ctx := context.WithValue(r.Context(), userIDKey, userID)
        next(w, r.WithContext(ctx))
    }
}

// GetUserIDFromContext lấy user_id từ context
func GetUserIDFromContext(ctx context.Context) string {
    if userID, ok := ctx.Value(userIDKey).(string); ok {
        return userID
    }
    return ""
}