package services

import (
	"consumer/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type RedisStatsRepo struct {
	rdb  *redis.Client
	pipe redis.Pipeliner
}

func NewRedisStatsRepo(cfg RedisConfig) *RedisStatsRepo {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
	})
	pipe := rdb.Pipeline()
	return &RedisStatsRepo{rdb, pipe}
}

func (r *RedisStatsRepo) GetStats(ctx context.Context, key string) (*models.Stats, error) {
	volumeKey := key + ":volume"
	volume, err := r.rdb.Get(ctx, volumeKey).Float64()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get volume data from redis for key %s", volumeKey)
	}

	txCountKey := key + ":tx_count"
	txCount, err := r.rdb.Get(ctx, txCountKey).Int()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get tx count data from redis for key %s", volumeKey)
	}

	return &models.Stats{
		Volume:  volume,
		TxCount: txCount,
	}, nil
}

func (r *RedisStatsRepo) UpsertStats(ctx context.Context, key string, value float64) error {
	now := time.Now().Unix()
	key5min := fmt.Sprintf("stats:%s:5min", key)
	key1h := fmt.Sprintf("stats:%s:1h", key)
	key24h := fmt.Sprintf("stats:%s:24h", key)

	r.pipe.IncrByFloat(ctx, key5min+":volume", value)
	r.pipe.Incr(ctx, key5min+":tx_count")
	r.pipe.ExpireAt(ctx, key5min, time.Unix(now+300, 0)) // Expire after 5 minutes

	r.pipe.IncrByFloat(ctx, key1h+":volume", value)
	r.pipe.Incr(ctx, key1h+":tx_count")
	r.pipe.ExpireAt(ctx, key1h, time.Unix(now+3600, 0)) // Expire after 1 hour

	r.pipe.IncrByFloat(ctx, key24h+":volume", value)
	r.pipe.Incr(ctx, key24h+":tx_count")
	r.pipe.ExpireAt(ctx, key24h, time.Unix(now+86400, 0)) // Expire after 24 hours

	_, err := r.pipe.Exec(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to upsert stats into Redis via its pipeline for key %s and value %v", key, value)
	}

	return nil
}
