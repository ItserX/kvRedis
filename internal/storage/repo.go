package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"kvManager/internal/pkg/log"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

func (repo *RedisRepository) marshalValue(value any) (string, error) {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("failed to marshal value: %v", err)
	}
	return string(jsonData), nil
}

func (repo *RedisRepository) unmarshalValue(data string) (any, error) {
	var result any
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal value: %v", err)
	}
	return result, nil
}

func (repo *RedisRepository) AddValue(key string, value any) error {
	log.Logger.Debugw("Adding value to Redis", "key", key)

	marshaledValue, err := repo.marshalValue(value)
	if err != nil {
		return err
	}

	_, err = repo.client.Set(context.Background(), key, marshaledValue, 0).Result()
	return err
}

func (repo *RedisRepository) GetValue(key string) ([]any, error) {
	log.Logger.Debugw("Get value from Redis", "key", key)

	data, err := repo.client.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			log.Logger.Debugw("Key not found in Redis")
			return nil, ErrKeyNotFound
		}
		log.Logger.Warnw("Redis operation failed", "error", err.Error())
		return nil, err
	}

	unmarshaled, err := repo.unmarshalValue(data)
	if err != nil {
		return nil, err
	}

	result := make([]any, 2)
	result[0] = key
	result[1] = unmarshaled

	return result, nil
}

func (repo *RedisRepository) UpdateValue(key string, value any) error {
	return repo.AddValue(key, value)
}

func (repo *RedisRepository) DeleteValue(key string) error {
	log.Logger.Debugw("Delete value from Redis", "key", key)
	_, err := repo.client.Del(context.Background(), key).Result()
	return err
}
