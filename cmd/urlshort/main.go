package main

import (
    "log"
    "net/http"
    "os"

    "urlshortener/internal/handlers"
    "urlshortener/internal/service"
    "urlshortener/internal/storage"

    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables from .env file
    godotenv.Load(".env")

    // Kết nối Redis
    store := storage.NewRedisStore()

    // Tạo services
    shortenerService := service.NewShortenerService(store)
    userService := service.NewUserService(store)  // Fix: Tạo đúng instance

    // Tạo handlers
    handler := handlers.NewHandler(shortenerService)
    authHandler := handlers.NewAuthHandler(userService, os.Getenv("JWT_SECRET_KEY"))  // Fix: Dùng userService đúng

    // Routes
    http.HandleFunc("/signup", authHandler.SignupHandler)
    http.HandleFunc("/login", authHandler.LoginHandler)
    http.HandleFunc("/shorten", authHandler.AuthMiddleware(handler.ShortenHandler))  // Fix: Thêm middleware để require auth
    http.HandleFunc("/", handler.RedirectHandler)  // Không cần auth

    log.Println("🚀 URL Shortener đang chạy tại http://localhost:8080")
    log.Println("Ví dụ: http://localhost:8080/shorten?url=https://google.com")

    log.Fatal(http.ListenAndServe(":8080", nil))
}