package service

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"time"
	"urlshortener/internal/models"
	"urlshortener/internal/storage"

	"github.com/redis/go-redis/v9"
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
	
	data, err := json.Marshal(url)
	if err != nil {
		return "", err
	}

	key := s.store.GetPrefix() + shortCode
	err = s.store.GetClient().Set(s.store.GetContext(), key, data, 0).Err()
	if err != nil {
		return "", err
	}

	return shortCode , nil
}

// GetLongURL lấy URL gốc từ short code
func (s *ShortenerService) GetLongURL(shortCode string) (string, bool) {
	key := s.store.GetPrefix() + shortCode
	data, err := s.store.GetClient().Get(s.store.GetContext(), key).Result()
	if err != nil {
		if  errors.Is(err, redis.Nil) {
			return "", false
		}
		return "", false
	}
	var url models.URL
	if err := json.Unmarshal([]byte(data), &url); err != nil {
		return "", false
	}
	url.Clicks += 1

	updateData, err := json.Marshal(url)
	if err != nil {
		log.Printf("❌ Lỗi khi cập nhật lượt click: %v", err)
		return "", false
	}
	err = s.store.GetClient().Set(s.store.GetContext(), key, updateData, 0).Err()
	if err != nil{
		return "", false
	}
	return url.LongURL, true
}

// generateShortCode tạo mã ngắn ngẫu nhiên
func generateShortCode() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:6]
}