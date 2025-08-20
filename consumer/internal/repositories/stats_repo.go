package repositories

import (
	"consumer/internal/models"
	"context"
)

type StatsRepo interface {
	GetStats(ctx context.Context, key string) (*models.Stats, error)
	UpsertStats(ctx context.Context, key string, value float64) error
}
