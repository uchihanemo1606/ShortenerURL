package storage

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
	prefix string // dùng để tổ chức key: "url:abc123"
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

// Save lưu dữ liệu short URL dưới dạng Hash
func (r *RedisStore) Save(shortCode, longURL string) {
	key := r.prefix + shortCode

	data := map[string]interface{}{
		"long":    longURL,
		"created": time.Now().Unix(),
		"clicks":  0,
	}

	// Lưu toàn bộ hash
	if err := r.client.HSet(r.ctx, key, data).Err(); err != nil {
		log.Printf("Redis lưu lỗi: %v", err)
	}

	// Optional: set thời hạn hết hạn (ví dụ 1 năm)
	// r.client.Expire(r.ctx, key, 365*24*time.Hour)
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