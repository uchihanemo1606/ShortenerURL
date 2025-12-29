package storage

import (
    "context"
    "log"
    "os"
    "github.com/redis/go-redis/v9"
)

type RedisStore struct {
    client *redis.Client
    ctx    context.Context
    prefix string
}


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

// GetClient trả về client Redis để service sử dụng
func (r *RedisStore) GetClient() *redis.Client {
    return r.client
}

// GetContext trả về context để service sử dụng
func (r *RedisStore) GetContext() context.Context {
    return r.ctx
}

// GetPrefix trả về prefix để service sử dụng
func (r *RedisStore) GetPrefix() string {
    return r.prefix
}