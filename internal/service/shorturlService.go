package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	// "hash"
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
func (s *ShortenerService) ShortenURL(longURL string, UserID string) (string, error) {

	if shortCode, exits := s.GetExitingShortCode(longURL); exits {
		return shortCode, nil
	}

	lockkey := s.store.GetPrefix() + "lock:" + fmt.Sprintf("%x", sha256.Sum256([]byte(longURL)))
	locked, err := s.store.GetClient().SetNX(s.store.GetContext(), lockkey, "1", 5*time.Second).Result()
	if err != nil {
		return "", err
	}

	if !locked {
		time.Sleep(100 * time.Millisecond)
		return s.ShortenURL(longURL, UserID)
	}
	defer s.store.GetClient().Del(s.store.GetContext(), lockkey)

	// Tạo short code ngẫu nhiên 6 ký tự
	shortCode := s.generateUniqueShortCode()

	if shortCode, exists := s.GetExitingShortCode(longURL); exists {
		return shortCode, nil
	}

	url := models.URL{
		ShortCode: shortCode,
		LongURL:   longURL,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(14 * 24 * time.Hour),
		UserID: UserID,
		Clicks:    0,
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

	return shortCode, nil
}

// GetLongURL lấy URL gốc từ short code
func (s *ShortenerService) GetLongURL(shortCode string) (string, bool) {
	key := s.store.GetPrefix() + shortCode
	data, err := s.store.GetClient().Get(s.store.GetContext(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
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
		log.Printf("Error when increasing clicks : %v", err)
		return "", false
	}
	err = s.store.GetClient().Set(s.store.GetContext(), key, updateData, 0).Err()
	if err != nil {
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

func (s *ShortenerService) GetExitingShortCode(longURL string) (string, bool) {
	hash := sha256.Sum256([]byte(longURL))
	key := s.store.GetPrefix() + "long:" + fmt.Sprintf("%x", hash)

	shortCode, err := s.store.GetClient().Get(s.store.GetContext(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", false
		}
		log.Printf("Error: check exit long url: %v", err)
		return "", false
	}
	return shortCode, true
}

func (s *ShortenerService) generateUniqueShortCode() string {
	for {
		shortCode := generateShortCode()
		keyShort := s.store.GetPrefix() + shortCode
		if _, err := s.store.GetClient().Get(s.store.GetContext(), keyShort).Result(); errors.Is(err, redis.Nil) {
			return shortCode
		}
	}
}

func (s *ShortenerService) GetAllShortURLs() ([]models.URL, error){
	pattern := s.store.GetPrefix() + "*"
	keys, err := s.store.GetClient().Keys(s.store.GetContext(), pattern).Result()
	if err != nil {
		return nil, err
	}

	var urls []models.URL
	for _, key := range keys {
		data, err := s.store.GetClient().Get(s.store.GetContext(), key).Result()
		if err != nil {
			log.Printf("Error retrieving URL for key %s: %v", key, err)
			continue
		}
		var url models.URL
		if err := json.Unmarshal([]byte(data), &url); err != nil {
			log.Printf("Error unmarshaling URL for key %s: %v", key, err)
			continue
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func (s *ShortenerService) ValidateLongURL(LongURL string) error {
	if LongURL == "" {
		return errors.New("Long URL is required")
	}
	if len(LongURL) > 2048 {
		return errors.New("Long URL exceeds maximum length of 2048 characters")
	}

	parsedURL, err := url.Parse(LongURL)
    if err != nil {
        return errors.New("URL not valid")
    }

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
        return errors.New("Only http and https URLs are supported")
    }

	hostname := strings.ToLower(parsedURL.Hostname())
    if hostname == "" || hostname == "localhost" || strings.HasPrefix(hostname, "127.") || strings.HasPrefix(hostname, "192.168.") || strings.HasPrefix(hostname, "10.") {
        return errors.New("Invalid hostname")
    }
	return nil

}