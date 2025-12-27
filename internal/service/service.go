package service

import (
	"crypto/rand"
	"encoding/base64"
	"urlshortener/internal/storage"
)

type ShortenerService struct {
	store *storage.RedisStore
}

func NewShortenerService(store *storage.RedisStore) *ShortenerService {
	return &ShortenerService{
		store: store,
	}
}

// ShortenURL tạo short code cho URL dài
func (s *ShortenerService) ShortenURL(longURL string) string {
	// Tạo short code ngẫu nhiên 6 ký tự
	shortCode := generateShortCode()

	// Lưu vào storage
	s.store.Save(shortCode, longURL)

	return shortCode
}

// GetLongURL lấy URL gốc từ short code
func (s *ShortenerService) GetLongURL(shortCode string) (string, bool) {
	return s.store.FindByShortCode(shortCode)
}

// generateShortCode tạo mã ngắn ngẫu nhiên
func generateShortCode() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:6]
}
