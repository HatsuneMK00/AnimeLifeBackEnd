package services

import (
	"AnimeLifeBackEnd/global"
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type AuthService interface {
	SaveCode(email string, code int) error
	GetCode(email string) (string, bool)
}

type authService struct{}

func (a authService) SaveCode(email string, code int) error {
	ctx := context.Background()
	err := global.RedisDB.Set(ctx, email, code, time.Hour).Err()
	return err
}

func (a authService) GetCode(email string) (string, bool) {
	ctx := context.Background()
	code, err := global.RedisDB.Get(ctx, email).Result()
	if err == redis.Nil {
		global.Logger.Errorf("no code for email or expired: %v", email)
		return "", false
	} else if err != nil {
		global.Logger.Errorf("get code from redis failed: %v", err)
		return "", false
	}
	return code, true
}
