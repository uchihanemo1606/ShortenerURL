package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"urlshortener/internal/service"
)

type Handler struct {
	service *service.ShortenerService
}

func NewHandler(service *service.ShortenerService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := GetUserIDFromContext(r.Context())
	if userID == ""{
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Lấy URL từ query parameter
	longURL := r.URL.Query().Get("url")
	if longURL == "" {
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}

	// Validate URL
	if _, err := url.Parse(longURL); err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// Tạo short URL
	shortCode, err := h.service.ShortenURL(longURL, userID)
	if err != nil {
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}
	shortURL := os.Getenv("BASE_URL") + shortCode

	// Trả về JSON response
	response := map[string]string{
		"short_url": shortURL,
		"long_url":  longURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RedirectHandler xử lý redirect từ short URL
func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {

	shortCode := r.URL.Path[1:]
	if shortCode == "" {
		http.NotFound(w, r)
		return
	}

	longURL, found := h.service.GetLongURL(shortCode)
	if !found {
		http.NotFound(w,r)
		return
	}
	
	http.Redirect(w, r, longURL, http.StatusFound)
}


func (h *Handler) GetAllShortURLsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    urls, err := h.service.GetAllShortURLs() 
    if err != nil {
        http.Error(w, "Failed to get URLs", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(urls)
}