package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	SetScheduleCache(ctx context.Context, studentID string, data []byte) error
	GetScheduleCache(ctx context.Context, studentID string) ([]byte, error)
}

type redisRepo struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepo{client: client}
}

func (r *redisRepo) SetScheduleCache(ctx context.Context, studentID string, data []byte) error {
	// Кэшируем расписание на 1 час
	return r.client.Set(ctx, "schedule:"+studentID, data, time.Hour).Err()
}

func (r *redisRepo) GetScheduleCache(ctx context.Context, studentID string) ([]byte, error) {
	return r.client.Get(ctx, "schedule:"+studentID).Bytes()
}
