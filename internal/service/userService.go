package service

import (
	"encoding/json"
	"errors"
	"time"
	"urlshortener/internal/models"
	"urlshortener/internal/storage"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	// "encoding/json"
)

type UserService struct {
	store *storage.RedisStore
}

func NewUserService(store *storage.RedisStore) *UserService {
	return &UserService{
		store: store,
	}
}


func (us *UserService) CreateUser(email, password string) (* models.User, error) {

	key := us.store.GetPrefix() + "user:" + email
	if _, err := us.store.GetClient().Get(us.store.GetContext(),key).Result(); err == nil {
		return nil, errors.New("user already exists")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		ID:  uuid.New().String(),
		Email: email,
		Password: string(hashPassword),
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	err = us.store.GetClient().Set(us.store.GetContext(), key, data, 0).Err()
	if err != nil {
		return nil, err
	}

	return user, nil
}