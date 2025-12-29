package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    "urlshortener/internal/service"
    "github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
    userService *service.UserService
    jwtSecret   []byte // Nên lấy từ env, ví dụ: os.Getenv("JWT_SECRET")
}

func NewAuthHandler(userService *service.UserService, jwtSecret string) *AuthHandler {
    return &AuthHandler{
        userService: userService,
        jwtSecret:   []byte(jwtSecret),
    }
}

// SignupHandler xử lý đăng ký
func (ah *AuthHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    user, err := ah.userService.CreateUser(req.Email, req.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "User created", "id": user.ID})
}

// LoginHandler xử lý đăng nhập
func (ah *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    user, err := ah.userService.AuthenticateUser(req.Email, req.Password)
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    // Tạo JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "exp":     time.Now().Add(time.Hour * 24).Unix(), // Expire sau 24h
    })
    tokenString, err := token.SignedString(ah.jwtSecret)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}