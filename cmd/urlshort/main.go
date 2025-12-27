package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"urlshortener/internal/handlers"
	"urlshortener/internal/service"
	"urlshortener/internal/storage"

)

func main() {
	// Load environment variables from .env file
	godotenv.Load(".env")

	// Kết nối Redis
	store := storage.NewRedisStore()

	// Tạo service và handler
	shortenerService := service.NewShortenerService(store)
	handler := handlers.NewHandler(shortenerService)

	// Routes
	http.HandleFunc("/shorten", handler.ShortenHandler)	
	http.HandleFunc("/", handler.RedirectHandler)

	log.Println("🚀 URL Shortener đang chạy tại http://localhost:8080")
	log.Println("Ví dụ: http://localhost:8080/shorten?url=https://google.com")

	log.Fatal(http.ListenAndServe(":8080", nil))
}