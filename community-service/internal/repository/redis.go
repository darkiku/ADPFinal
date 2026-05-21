package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	AcquireLock(ctx context.Context, spaceID, timeSlot string) (bool, error)
	ReleaseLock(ctx context.Context, spaceID, timeSlot string) error
}

type redisRepo struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepo{client: client}
}

// AcquireLock пытается захватить блокировку с помощью SetNX (Set if Not eXists)
func (r *redisRepo) AcquireLock(ctx context.Context, spaceID, timeSlot string) (bool, error) {
	key := fmt.Sprintf("lock:coworking:%s:%s", spaceID, timeSlot)

	// Блокируем на 5 минут (время на завершение транзакции)
	locked, err := r.client.SetNX(ctx, key, "locked", 5*time.Minute).Result()
	if err != nil {
		return false, err
	}
	return locked, nil
}

func (r *redisRepo) ReleaseLock(ctx context.Context, spaceID, timeSlot string) error {
	key := fmt.Sprintf("lock:coworking:%s:%s", spaceID, timeSlot)
	return r.client.Del(ctx, key).Err()
}
