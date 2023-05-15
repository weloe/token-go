package redis_adapter

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/weloe/token-go/token-go/model"
	"github.com/weloe/token-go/token-go/persist"
	"time"
)

var _ persist.Adapter = (*RedisAdapter)(nil)

type RedisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(addr string, password string, db int) (*RedisAdapter, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return &RedisAdapter{client: client}, nil
}

func (r *RedisAdapter) GetStr(key string) string {
	res, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		return ""
	}
	return res
}

func (r *RedisAdapter) SetStr(key string, value string, timeout int64) error {
	err := r.client.Set(context.Background(), key, value, time.Duration(timeout)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisAdapter) UpdateStr(key string, value string) error {
	err := r.client.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisAdapter) DeleteStr(key string) error {
	err := r.client.Del(context.Background(), key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisAdapter) GetStrTimeout(key string) int64 {
	duration, err := r.client.TTL(context.Background(), key).Result()
	if err != nil {
		return -1
	}
	return int64(duration.Seconds())
}

func (r *RedisAdapter) UpdateStrTimeout(key string, timeout int64) error {
	var duration time.Duration
	if timeout < 0 {
		duration = -1
	} else {
		duration = time.Duration(timeout) * time.Second
	}
	err := r.client.Expire(context.Background(), key, duration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisAdapter) Get(key string) interface{} {
	res, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil
	}
	s := &model.Session{}
	err = json.Unmarshal([]byte(res), s)
	if err != nil {
		return nil
	}
	return s
}

func (r *RedisAdapter) Set(key string, value interface{}, timeout int64) error {
	err := r.client.Set(context.Background(), key, value, time.Duration(timeout)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisAdapter) Update(key string, value interface{}) error {
	err := r.client.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisAdapter) Delete(key string) error {
	err := r.client.Del(context.Background(), key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisAdapter) GetTimeout(key string) int64 {
	duration, err := r.client.TTL(context.Background(), key).Result()
	if err != nil {
		return -1
	}
	return int64(duration.Seconds())
}

func (r *RedisAdapter) UpdateTimeout(key string, timeout int64) error {
	var duration time.Duration
	if timeout < 0 {
		duration = -1
	} else {
		duration = time.Duration(timeout) * time.Second
	}
	err := r.client.Expire(context.Background(), key, duration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisAdapter) DeleteBatchFilteredKey(filterKeyPrefix string) error {
	var cursor uint64 = 0
	for {
		keys, cursor, err := r.client.Scan(context.Background(), cursor, filterKeyPrefix+"*", 100).Result()
		if err != nil {
			return err
		}

		if len(keys) == 0 && cursor == 0 {
			break
		}

		// use pip delete batch
		pipe := r.client.Pipeline()

		for _, key := range keys {
			pipe.Del(context.Background(), key)
		}

		_, err = pipe.Exec(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}
