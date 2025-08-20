package services

import (
	"consumer/internal/models"
	"consumer/internal/utils"
	"context"
	"fmt"
	"strconv"
	"time"

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

// Get stats via a key "stats:ETH:5min" with consideration to the window bucket
// Each bucket key is a postfix for the original key
// O(k) reads, where k = buckets in window (5, 12, or 24 max)
func (r *RedisStatsRepo) GetStats(ctx context.Context, key string) (*models.Stats, error) {
	var volume float64
	var txCount int64

	bucketKeys := getWindowBuckets(key)
	for _, bucketKey := range bucketKeys {
		volResult := r.rdb.Get(ctx, bucketKey+":volume")
		if volResult.Err() == nil {
			if val, err := volResult.Float64(); err == nil {
				volume += val
			}
		}

		countResult := r.rdb.Get(ctx, bucketKey+":tx_count")
		if countResult.Err() == nil {
			if val, err := countResult.Int64(); err == nil {
				txCount += val
			}
		}
	}

	return &models.Stats{
		Volume:  volume,
		TxCount: txCount,
	}, nil
}

// Method to aggregate stats data
// - O(1) writes: 6 Redis operations per swap
// - Fixed memory
// - Automatic cleanup: Redis TTL handles expiration
func (r *RedisStatsRepo) UpsertStats(ctx context.Context, key string, value float64) error {
	now := time.Now()

	// 5min buckets (60 seconds each)
	bucket5min := now.Truncate(time.Minute).Unix()
	key5min := fmt.Sprintf("stats:%s:5min:%d", key, bucket5min)

	// 1h buckets (5 minutes each)
	bucket1h := now.Truncate(5 * time.Minute).Unix()
	key1h := fmt.Sprintf("stats:%s:1h:%d", key, bucket1h)

	// 24h buckets (1 hour each)
	bucket24h := now.Truncate(time.Hour).Unix()
	key24h := fmt.Sprintf("stats:%s:24h:%d", key, bucket24h)

	r.pipe.IncrByFloat(ctx, key5min+":volume", value)
	r.pipe.Incr(ctx, key5min+":tx_count")
	r.pipe.Expire(ctx, key5min+":volume", 5*time.Minute)
	r.pipe.Expire(ctx, key5min+":tx_count", 5*time.Minute)

	r.pipe.IncrByFloat(ctx, key1h+":volume", value)
	r.pipe.Incr(ctx, key1h+":tx_count")
	r.pipe.Expire(ctx, key1h+":volume", 60*time.Minute)
	r.pipe.Expire(ctx, key1h+":tx_count", 60*time.Minute)

	r.pipe.IncrByFloat(ctx, key24h+":volume", value)
	r.pipe.Incr(ctx, key24h+":tx_count")
	r.pipe.Expire(ctx, key24h+":volume", 24*time.Hour)
	r.pipe.Expire(ctx, key24h+":tx_count", 24*time.Hour)

	_, err := r.pipe.Exec(ctx)
	return err
}

// Get the bucket keys based on the current time and provided original key
// - 5min window divided into 5 buckets of 1 minute each
// - 1h window divided into 12 buckets of 5 minutes each
// - 24h window - 24 buckets of 1 hour each
func getWindowBuckets(originalKey string) []string {
	now := time.Now()
	var buckets []string
	if utils.Contains(originalKey, ":5min") {
		for i := 0; i < 5; i++ {
			bucket := now.Add(-time.Duration(i) * time.Minute).Truncate(time.Minute).Unix()
			buckets = append(buckets, buildBucketKey(originalKey, bucket))
		}
	}
	if utils.Contains(originalKey, ":1h") {
		for i := 0; i < 12; i++ {
			bucket := now.Add(-time.Duration(i) * 5 * time.Minute).Truncate(5 * time.Minute).Unix()
			buckets = append(buckets, buildBucketKey(originalKey, bucket))
		}
	}
	if utils.Contains(originalKey, ":24h") {
		for i := 0; i < 24; i++ {
			bucket := now.Add(-time.Duration(i) * time.Hour).Truncate(time.Hour).Unix()
			buckets = append(buckets, buildBucketKey(originalKey, bucket))
		}
	}
	return buckets
}

// Convert the original key to the bucket key
// For example, "stats:ETH:5min" + bucket -> "stats:ETH:5min:1735689600"
func buildBucketKey(originalKey string, bucket int64) string {
	return originalKey + ":" + strconv.FormatInt(bucket, 10)
}
