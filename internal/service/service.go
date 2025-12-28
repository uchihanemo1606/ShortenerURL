package service

import (
	"crypto/rand"
	"encoding/base64"
	"time"
	"urlshortener/internal/models"
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
func (s *ShortenerService) ShortenURL(longURL string) (string, error) {
	// Tạo short code ngẫu nhiên 6 ký tự
	shortCode := generateShortCode()

	url := models.URL{
		ShortCode : shortCode,
		LongURL : longURL,
		CreatedAt : time.Now(),
		Clicks : 0,
	}
	if err := s.store.Save(url); err != nil {
        return "", err
    }
	return shortCode , nil
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