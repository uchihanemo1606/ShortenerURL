package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
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

// ShortenHandler xử lý yêu cầu tạo short URL
func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
	shortCode, err := h.service.ShortenURL(longURL)
	if err != nil {
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}
	shortURL := "http://localhost:8080/" + shortCode

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
	// Lấy short code từ path
	shortCode := r.URL.Path[1:]
	if shortCode == "" {
		http.NotFound(w, r)
		return
	}

	// Tìm URL gốc
	longURL, found := h.service.GetLongURL(shortCode)
	if !found {
		http.NotFound(w,r)
		return
	}
	

	// Redirect
	http.Redirect(w, r, longURL, http.StatusFound)
}