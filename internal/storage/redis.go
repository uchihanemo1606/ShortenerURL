package storage

import (
	"context"
	"log"
	"os"
	"github.com/redis/go-redis/v9"
	"urlshortener/internal/models"
	"encoding/json"
)

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
	prefix string 
}

// NewRedisStore khởi tạo kết nối với Upstash Redis
func NewRedisStore() *RedisStore {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("❌ REDIS_URL không được set! Hãy set biến môi trường.")
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("❌ Sai định dạng Redis URL: %v", err)
	}

	client := redis.NewClient(opt)
	ctx := context.Background()

	// Test kết nối ngay khi khởi tạo
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("❌ Không kết nối được với Redis: %v", err)
	}

	log.Println("✅ Kết nối Upstash Redis thành công!")

	return &RedisStore{
		client: client,
		ctx:    ctx,
		prefix: "url:",
	}
}

func (r *RedisStore) Save(url models.URL) error {
	key := r.prefix + url.ShortCode

	// Chuyển struct thành JSON
	data, err := json.Marshal(url)
	if err != nil {
		log.Printf("Lỗi marshal JSON: %v", err)
		return err
	}

	// Lưu JSON vào Redis
	if err := r.client.Set(r.ctx, key, data, 0).Err(); err != nil {
		log.Printf("Redis lưu lỗi: %v", err)
		return err
	}

	return nil
}

// FindByShortCode lấy URL gốc và tăng lượt click
func (r *RedisStore) FindByShortCode(shortCode string) (string, bool) {
	key := r.prefix + shortCode

	longURL, err := r.client.HGet(r.ctx, key, "long").Result()
	if err == redis.Nil {
		return "", false // không tìm thấy
	}
	if err != nil {
		log.Printf("Redis đọc lỗi: %v", err)
		return "", false
	}

	// Tăng lượt click (atomic - an toàn khi nhiều người truy cập cùng lúc)
	r.client.HIncrBy(r.ctx, key, "clicks", 1)

	return longURL, true
}